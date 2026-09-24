package poryscriptz

import (
	"fmt"
	"strconv"
	"strings"

	gts "github.com/odvcencio/gotreesitter"
	"github.com/odvcencio/gotreesitter/taproot"
)

// Event patterns expand to existing HGSS script.inc commands. They do not
// allocate opcodes or change the argument widths of the underlying commands.
type patternArg string

const (
	messageID    patternArg = "message id"
	trainerID    patternArg = "trainer id"
	itemID       patternArg = "item id"
	giftQuantity patternArg = "gift quantity"
	flagID       patternArg = "FLAG_ identifier"
	varID        patternArg = "VAR_ identifier"
	mapID        patternArg = "map id"
	actorID      patternArg = "actor id"
	labelID      patternArg = "label identifier"
	scalar       patternArg = "integer or symbol"
	compareValue patternArg = "literal or constant value (not VAR_)"
	direction    patternArg = "direction"
)

var patternArgs = map[string][]patternArg{
	"message":               {messageID},
	"say":                   {messageID},
	"trainer_battle_simple": {trainerID},
	"give_item":             {itemID, giftQuantity, labelID},
	"move_actor":            {actorID, labelID},
	"move_actor_and_wait":   {actorID, labelID},
	"flag_set":              {flagID},
	"flag_clear":            {flagID},
	"var_set":               {varID, compareValue},
	"var_copy":              {varID, varID},
	"warp_to":               {mapID, scalar, scalar, scalar, direction},
	"if_flag_set":           {flagID, labelID},
	"if_flag_unset":         {flagID, labelID},
	"if_var_eq":             {varID, compareValue, labelID},
	"if_var_ne":             {varID, compareValue, labelID},
}

func resolvePattern(w *taproot.Walker, name string, argList *gts.Node, args []string, line, col int, targetName string) (Stmt, []Diag, bool) {
	kinds, ok := patternArgs[name]
	if !ok || targetName != "HGSS" {
		return nil, nil, false
	}
	if len(args) != len(kinds) {
		return nil, []Diag{{Line: line, Col: col, Msg: fmt.Sprintf("pattern %s expects %d args, got %d", name, len(kinds), len(args))}}, true
	}
	var diags []Diag
	for i, kind := range kinds {
		n := argList.NamedChild(i)
		argLine, argCol := w.Pos(n)
		if !validPatternArg(w.Type(n), args[i], kind) {
			diags = append(diags, Diag{Line: argLine, Col: argCol,
				Msg: fmt.Sprintf("pattern %s arg %d requires %s", name, i+1, kind)})
			continue
		}
		if w.Type(n) == "int_literal" {
			max := patternLiteralMax(kind)
			value, err := strconv.ParseUint(args[i], 0, 64)
			if err != nil || value > max {
				diags = append(diags, Diag{Line: argLine, Col: argCol,
					Msg: fmt.Sprintf("pattern %s arg %d numeric value must be <= %d", name, i+1, max)})
			}
		}
	}
	if len(diags) > 0 {
		return nil, diags, true
	}
	if targetIndex, _, hasTarget := patternTarget(name); hasTarget {
		line, col = w.Pos(argList.NamedChild(targetIndex))
	}
	return PatternCall{Name: name, Args: args, Line: line, Col: col}, nil, true
}

func patternLiteralMax(kind patternArg) uint64 {
	switch kind {
	case messageID:
		return 0xFF // npc_msg writes one byte
	case itemID, giftQuantity:
		return 0x3FFF // item_vars copies a VAR_* when the value is >= 0x4000
	default:
		return 0xFFFF // all other numeric pattern fields are .short
	}
}

func validPatternArg(nodeType, text string, kind patternArg) bool {
	ident := nodeType == "identifier"
	integer := nodeType == "int_literal"
	switch kind {
	case flagID:
		return ident && strings.HasPrefix(text, "FLAG_")
	case varID:
		return ident && strings.HasPrefix(text, "VAR_")
	case itemID:
		return integer || ident && (strings.HasPrefix(text, "ITEM_") || strings.HasPrefix(text, "TM_") || strings.HasPrefix(text, "VAR_"))
	case giftQuantity:
		return integer || ident && !strings.HasPrefix(text, "FLAG_")
	case trainerID:
		return integer || ident && strings.HasPrefix(text, "TRAINER_")
	case messageID:
		return integer || ident && (strings.HasPrefix(text, "msg_") || strings.HasPrefix(text, "MSG_"))
	case mapID:
		return integer || ident && strings.HasPrefix(text, "MAP_")
	case actorID:
		return integer || ident && (strings.HasPrefix(text, "obj_") || strings.HasPrefix(text, "OBJ_") || strings.HasPrefix(text, "VAR_"))
	case direction:
		return integer || ident && strings.HasPrefix(text, "DIR_")
	case labelID:
		return ident
	case scalar:
		return integer || ident
	case compareValue:
		return integer || ident && !strings.HasPrefix(text, "VAR_")
	default:
		return false
	}
}

// patternTarget identifies arguments that must resolve to a local section.
func patternTarget(name string) (index int, movement bool, ok bool) {
	switch name {
	case "give_item", "if_var_eq", "if_var_ne":
		return 2, false, true
	case "if_flag_set", "if_flag_unset":
		return 1, false, true
	case "move_actor", "move_actor_and_wait":
		return 1, true, true
	default:
		return 0, false, false
	}
}

func validatePatternTargets(prog *Program) []Diag {
	branchLabels := make(map[string]bool)
	movementLabels := make(map[string]bool)
	for _, section := range prog.Sections {
		if section.Align {
			continue
		}
		if section.Movement {
			movementLabels[section.Name] = true
		} else {
			branchLabels[section.Name] = true
		}
	}
	var diags []Diag
	var visit func([]Stmt)
	visit = func(stmts []Stmt) {
		for _, stmt := range stmts {
			switch st := stmt.(type) {
			case PatternCall:
				index, movement, ok := patternTarget(st.Name)
				if !ok {
					continue
				}
				label := st.Args[index]
				if movement && !movementLabels[label] {
					diags = append(diags, Diag{Line: st.Line, Col: st.Col,
						Msg: fmt.Sprintf("pattern %s requires a defined movement section %q", st.Name, label)})
				} else if !movement && !branchLabels[label] {
					diags = append(diags, Diag{Line: st.Line, Col: st.Col,
						Msg: fmt.Sprintf("pattern %s requires a defined branch label %q", st.Name, label)})
				}
			case If:
				visit(st.Body)
				visit(st.Else)
			case While:
				visit(st.Body)
			case Switch:
				for _, arm := range st.Cases {
					visit(arm.Body)
				}
				visit(st.Default)
			}
		}
	}
	for _, section := range prog.Sections {
		visit(section.Stmts)
	}
	return diags
}

func patternInstr(macro string, args ...string) Instr {
	return Instr{Opcode: -1, Macro: macro, Args: args}
}

// expandPattern is total for PatternCall values produced by resolvePattern.
func expandPattern(name string, args []string) []Instr {
	switch name {
	case "message":
		return []Instr{patternInstr("npc_msg", args[0])}
	case "say":
		return []Instr{patternInstr("npc_msg", args[0]), patternInstr("closemsg")}
	case "trainer_battle_simple":
		return []Instr{patternInstr("trainer_battle", args[0], "0", "0", "0")}
	case "give_item":
		return []Instr{patternInstr("goto_if_no_item_space", args[0], args[1], args[2]), patternInstr("callstd", "std_give_item_verbose")}
	case "move_actor":
		return []Instr{patternInstr("apply_movement", args[0], args[1])}
	case "move_actor_and_wait":
		return []Instr{patternInstr("apply_movement", args[0], args[1]), patternInstr("wait_movement")}
	case "flag_set":
		return []Instr{patternInstr("setflag", args[0])}
	case "flag_clear":
		return []Instr{patternInstr("clearflag", args[0])}
	case "var_set":
		return []Instr{patternInstr("setvar", args[0], args[1])}
	case "var_copy":
		return []Instr{patternInstr("copyvar", args[0], args[1])}
	case "warp_to":
		return []Instr{patternInstr("warp", args...)}
	case "if_flag_set":
		return []Instr{patternInstr("goto_if_set", args...)}
	case "if_flag_unset":
		return []Instr{patternInstr("goto_if_unset", args...)}
	case "if_var_eq":
		return []Instr{patternInstr("compare_var_to_value", args[0], args[1]), patternInstr("goto_if_eq", args[2])}
	case "if_var_ne":
		return []Instr{patternInstr("compare_var_to_value", args[0], args[1]), patternInstr("goto_if_ne", args[2])}
	default:
		panic("unknown resolved pattern: " + name)
	}
}
