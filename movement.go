package poryscriptz

import (
	"fmt"
	"strconv"

	gts "github.com/odvcencio/gotreesitter"
	"github.com/odvcencio/gotreesitter/taproot"
	"github.com/odvcencio/poryscriptz/target"
	"github.com/odvcencio/poryscriptz/vocab"
)

// movementArity lists a bounded set of HGSS script.inc movement macros.
// A range of 0..1 means the macro accepts an optional repeat count.
func movementArity(name string) (min, max int, ok bool) {
	switch name {
	case "EndMovement":
		return 0, 0, true
	case "Delay16", "EmoteExclamationMark", "SetVisible":
		return 0, 1, true
	}
	for _, dir := range []string{"North", "South", "East", "West"} {
		for _, stem := range []string{
			"WalkNormal", "WalkFast", "WalkSlightlyFast",
			"WalkOnSpotNormal", "WalkOnSpotFast", "JumpOnSpotFast",
		} {
			if name == stem+dir {
				return 0, 1, true
			}
		}
		if name == "JumpFar"+dir {
			return 0, 1, true
		}
	}
	return 0, 0, false
}

func resolveMovement(w *taproot.Walker, n *gts.Node, tgt target.GameTarget) (*Script, []Diag) {
	s := &Script{Name: w.Text(w.Field(n, "name"))}
	var diags []Diag
	block := w.ChildByType(n, "block")
	if block == nil {
		return s, diags
	}
	list := w.ChildByType(block, "statement_list")
	if list == nil {
		line, col := w.Pos(n)
		return s, []Diag{{Line: line, Col: col, Msg: "movement requires EndMovement()"}}
	}
	ended := false
	for i := 0; i < list.NamedChildCount(); i++ {
		stmt := list.NamedChild(i)
		call := w.ChildByType(stmt, "call_expression")
		if w.Type(stmt) != "expression_statement" || call == nil {
			line, col := w.Pos(stmt)
			diags = append(diags, Diag{Line: line, Col: col, Msg: "movement accepts movement macro calls only"})
			continue
		}
		nameNode := w.Field(call, "function")
		name := w.Text(nameNode)
		line, col := w.Pos(nameNode)
		if ended {
			diags = append(diags, Diag{Line: line, Col: col, Msg: "EndMovement must be the last movement command"})
			continue
		}
		min, max, ok := movementArity(name)
		if me, exists := tgt.Macros().ByName(name); exists && me.Movement {
			min, max, ok = me.MinArgs, len(me.Args), true
		}
		if !ok {
			diags = append(diags, Diag{Line: line, Col: col, Msg: fmt.Sprintf("unknown movement macro %q", name)})
			continue
		}
		argsNode := w.Field(call, "arguments")
		count := 0
		if argsNode != nil {
			count = argsNode.NamedChildCount()
		}
		if count < min || count > max {
			diags = append(diags, Diag{Line: line, Col: col, Msg: fmt.Sprintf("movement %s expects %d to %d args, got %d", name, min, max, count)})
			continue
		}
		var args []string
		if count == 1 {
			arg := argsNode.NamedChild(0)
			if w.Type(arg) != "int_literal" {
				argLine, argCol := w.Pos(arg)
				diags = append(diags, Diag{Line: argLine, Col: argCol, Msg: "movement repeat count must be an integer literal"})
				continue
			}
			text := w.Text(arg)
			if _, err := strconv.ParseUint(text, 0, 16); err != nil {
				argLine, argCol := w.Pos(arg)
				diags = append(diags, Diag{Line: argLine, Col: argCol, Msg: "movement repeat count must fit a 16-bit field"})
				continue
			}
			args = []string{text}
		}
		s.Stmts = append(s.Stmts, MacroCall{Macro: &vocab.MacroEntry{Name: name}, Args: args})
		if name == "EndMovement" {
			ended = true
		}
	}
	if !ended {
		line, col := w.Pos(n)
		diags = append(diags, Diag{Line: line, Col: col, Msg: "movement requires EndMovement()"})
	}
	return s, diags
}
