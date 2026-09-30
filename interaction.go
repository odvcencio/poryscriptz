package poryscriptz

import (
	"fmt"
	"strconv"
	"strings"

	gts "github.com/odvcencio/gotreesitter"
	"github.com/odvcencio/gotreesitter/taproot"
	"github.com/odvcencio/poryscriptz/target"
	"github.com/odvcencio/poryscriptz/vocab"
)

// An interaction is an NPC talk entry point. Its restricted body cannot jump
// out of the lock/message lifecycle. Raw scripts remain the escape hatch.
func resolveInteraction(w *taproot.Walker, n *gts.Node, tgt target.GameTarget) (*Script, []Diag) {
	s := &Script{Name: w.Text(w.Field(n, "name"))}
	if tgt.Name() != "HGSS" {
		line, col := w.Pos(n)
		return s, []Diag{{Line: line, Col: col, Msg: "interaction is supported only by the HGSS target"}}
	}
	line, col := w.Pos(n)
	if ds := requireHGSSCommands(tgt, []string{"lockall", "faceplayer", "releaseall", "end"}, line, col); len(ds) > 0 {
		return s, ds
	}
	block := w.ChildByType(n, "block")
	if diags := validateInteractionBody(w, block); len(diags) > 0 {
		return s, diags
	}
	resolved, diags := resolveScript(w, n, tgt)
	// Generated lifecycle commands are ordinary HGSS macros, never new opcodes.
	s.Stmts = append(s.Stmts, MacroCall{Macro: &vocab.MacroEntry{Name: "lockall"}}, MacroCall{Macro: &vocab.MacroEntry{Name: "faceplayer"}})
	s.Stmts = append(s.Stmts, resolved.Stmts...)
	s.Stmts = append(s.Stmts, MacroCall{Macro: &vocab.MacroEntry{Name: "releaseall"}}, MacroCall{Macro: &vocab.MacroEntry{Name: "end"}})
	return s, diags
}

// Validate the CST so rejected operations retain their exact source locations.
// Only helpers that always fall through and close their own message are allowed.
func validateInteractionBody(w *taproot.Walker, block *gts.Node) []Diag {
	var diags []Diag
	var visit func(*gts.Node)
	visit = func(n *gts.Node) {
		if n == nil {
			return
		}
		switch w.Type(n) {
		case "block", "statement_list":
			for i := 0; i < n.NamedChildCount(); i++ {
				visit(n.NamedChild(i))
			}
		case "comment":
		case "if_statement":
			// resolveIf checks the condition and rejects inherited Go initializers.
			diags = append(diags, validateInteractionCondition(w, w.Field(n, "condition"))...)
			visit(w.Field(n, "consequence"))
			visit(w.Field(n, "alternative"))
		case "expression_statement":
			call := w.ChildByType(n, "call_expression")
			name := w.Text(w.Field(call, "function"))
			switch name {
			case "dialogue", "flag_set", "flag_clear", "var_set", "var_copy":
				return
			}
			line, col := w.Pos(n)
			diags = append(diags, Diag{Line: line, Col: col, Msg: fmt.Sprintf("interaction does not allow %q; use dialogue, checked flag/variable helpers, and if/else, or a raw script for manual control flow", name)})
		default:
			line, col := w.Pos(n)
			diags = append(diags, Diag{Line: line, Col: col, Msg: "interaction does not allow " + w.Type(n) + "; use if/else or a raw script for manual control flow"})
		}
	}
	visit(block)
	return diags
}

func validateInteractionCondition(w *taproot.Walker, condition *gts.Node) []Diag {
	if condition == nil {
		return nil
	}
	var call *gts.Node
	var kind patternArg
	switch w.Type(condition) {
	case "call_expression":
		if w.Text(w.Field(condition, "function")) == "flag" {
			call, kind = condition, flagID
		}
	case "binary_expression":
		left := w.Field(condition, "left")
		if w.Type(left) == "call_expression" && w.Text(w.Field(left, "function")) == "var" {
			call, kind = left, varID
		}
	}
	if call == nil {
		return nil
	} // resolveIf rejects unsupported condition shapes.
	args := w.Field(call, "arguments")
	if args == nil || args.NamedChildCount() != 1 {
		return nil
	} // resolveIf reports arity.
	arg := args.NamedChild(0)
	line, col := w.Pos(arg)
	if !validPatternArg(w.Type(arg), w.Text(arg), kind) {
		return []Diag{{Line: line, Col: col, Msg: "interaction condition requires " + string(kind)}}
	}
	if kind == varID {
		right := w.Field(condition, "right")
		if w.Type(right) == "int_literal" {
			value, err := strconv.ParseUint(w.Text(right), 0, 64)
			if err != nil || value > 0xFFFF {
				line, col = w.Pos(right)
				return []Diag{{Line: line, Col: col, Msg: "interaction variable comparison must fit 16 bits (<= 65535)"}}
			}
		}
	}
	return nil
}

// Reuse dialogue's checked message reference rules without exposing ask as a
// statement pattern: its result belongs to the enclosing conditional.
func resolveAsk(w *taproot.Walker, call *gts.Node, tgt target.GameTarget) (string, []Diag) {
	line, col := w.Pos(call)
	if tgt.Name() != "HGSS" {
		return "", []Diag{{Line: line, Col: col, Msg: "ask is supported only by the HGSS target"}}
	}
	if ds := requireHGSSCommands(tgt, []string{"npc_msg", "yesno", "closemsg", "compare_var_to_value", "goto_if_ne"}, line, col); len(ds) > 0 {
		return "", ds
	}
	argList := w.Field(call, "arguments")
	var args []string
	for i := 0; argList != nil && i < argList.NamedChildCount(); i++ {
		args = append(args, w.Text(argList.NamedChild(i)))
	}
	_, diags, _ := resolvePattern(w, "dialogue", argList, args, line, col, tgt.Name())
	for i := range diags {
		// Explain the spelling the author actually used.
		diags[i].Msg = strings.Replace(diags[i].Msg, "pattern dialogue", "predicate ask", 1)
	}
	if len(diags) > 0 {
		return "", diags
	}
	return args[0], nil
}

func requireHGSSCommands(tgt target.GameTarget, names []string, line, col int) []Diag {
	for _, name := range names {
		if _, ok := tgt.Vocabulary().ByName(name); ok {
			continue
		}
		if _, ok := tgt.Macros().ByName(name); ok {
			continue
		}
		return []Diag{{Line: line, Col: col, Msg: fmt.Sprintf("HGSS event helper requires command %q in the selected vocabulary", name)}}
	}
	return nil
}
