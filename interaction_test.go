package poryscriptz

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/odvcencio/poryscriptz/target"
	"github.com/odvcencio/poryscriptz/vocab"
)

func compileFlow(t *testing.T, src string, tgt target.GameTarget) string {
	t.Helper()
	asm, diags, err := Compile([]byte(src), tgt)
	if err != nil || len(diags) != 0 {
		t.Fatalf("compile: err=%v diags=%v", err, diags)
	}
	return asm
}

func TestInteractionGreetingMatchesRawScript(t *testing.T) {
	managed := "package m\ninteraction Greeting { dialogue(MSG_GREETING) }\nalign4()\n"
	raw := "package m\nscript Greeting { lockall(); faceplayer(); npc_msg(MSG_GREETING); WaitButton(); closemsg(); releaseall(); end() }\nalign4()\n"
	bundled, err := target.HGSS("")
	if err != nil {
		t.Fatal(err)
	}
	for _, tgt := range []target.GameTarget{hgssT(t), bundled} {
		got, want := compileFlow(t, managed, tgt), compileFlow(t, raw, tgt)
		if got != want {
			t.Fatalf("managed greeting differs from raw script:\n%s\nwant:\n%s", got, want)
		}
		if again := compileFlow(t, managed, tgt); again != got {
			t.Fatal("nondeterministic output")
		}
	}
	root, w, err := Parse([]byte(managed))
	if err != nil || firstScriptName(w, root) != "Greeting" {
		t.Fatalf("interaction entry parse: %v", err)
	}
}

func TestDialogueAddsWaitWithoutChangingSay(t *testing.T) {
	tgt := hgssT(t)
	got := compileFlow(t, "package m\nscript S { dialogue(255); End() }", tgt)
	if !strings.Contains(got, "\tnpc_msg 255\n\tWaitButton\n\tclosemsg\n") {
		t.Fatalf("dialogue lifecycle: %s", got)
	}
	say := compileFlow(t, "package m\nscript S { say(1); End() }", tgt)
	if strings.Contains(say, "WaitButton") || !strings.Contains(say, "\tnpc_msg 1\n\tclosemsg\n") {
		t.Fatalf("say changed behavior: %s", say)
	}
}

func TestManagedBodyRejectsEscapesAtTheirSource(t *testing.T) {
	for _, operation := range []string{
		"End()", "AsmReturn()", "AsmGoto(_Away)", "GoTo(_Away)", "Call(_Away)",
		"LockAll()", "ReleaseAll()", "NPCMsg(1)", "say(1)", "warp_to(MAP_X, 0, 0, 0, DIR_NORTH)",
		"give_item(ITEM_POTION, 1, _Away)", "for flag(FLAG_X) {}", "return", "VAR_X = 1",
		"switch var(VAR_X) { case 0: dialogue(1) }",
		"1 + 2", "{ dialogue(1) }", "defer dialogue(1)",
	} {
		t.Run(operation, func(t *testing.T) {
			asm, diags, err := Compile([]byte("package m\ninteraction S {\n    "+operation+"\n}\nlabel _Away { End() }"), hgssT(t))
			if err != nil || asm != "" || len(diags) != 1 || diags[0].Line != 3 || diags[0].Col != 5 || !(strings.Contains(diags[0].Msg, "interaction does not allow") || strings.Contains(diags[0].Msg, "unsupported statement")) {
				t.Fatalf("unsafe operation escaped: err=%v asm=%q diags=%v", err, asm, diags)
			}
		})
	}
	asm, ds, err := Compile([]byte("package m\ninteraction S { if ask(0) { dialogue(1) } else if flag(FLAG_X) { End() } }"), hgssT(t))
	if err != nil || asm != "" || len(ds) != 1 || !strings.Contains(ds[0].Msg, `"End"`) {
		t.Fatalf("nested exit escaped: %v %v", err, ds)
	}
}

type testFlowTarget struct {
	target.GameTarget
	name          string
	withoutMacros bool
}

func (t testFlowTarget) Name() string { return t.name }
func (t testFlowTarget) Macros() *vocab.MacroTable {
	if t.withoutMacros {
		return nil
	}
	return t.GameTarget.Macros()
}

func TestFlowRequiresActualTargetCapabilities(t *testing.T) {
	base := hgssT(t)
	for _, tc := range []struct {
		source string
		tgt    target.GameTarget
		want   string
	}{
		{"interaction S { dialogue(1) }", testFlowTarget{GameTarget: base, name: "Other"}, "only by the HGSS target"},
		{"script S { if ask(1) {} }", testFlowTarget{GameTarget: base, name: "Other"}, "only by the HGSS target"},
		{"script S { dialogue(1) }", testFlowTarget{GameTarget: base, name: "HGSS", withoutMacros: true}, `requires command "WaitButton"`},
	} {
		asm, ds, err := Compile([]byte("package m\n"+tc.source), tc.tgt)
		if err != nil || asm != "" || len(ds) != 1 || !strings.Contains(ds[0].Msg, tc.want) || ds[0].Line != 2 {
			t.Fatalf("err=%v asm=%q diags=%v", err, asm, ds)
		}
	}
	for _, source := range []string{"if flag(VAR_X) {}", "if flag(1) {}", "if var(FLAG_X) == 0 {}", "if var(VAR_X) == 65536 {}"} {
		asm, ds, err := Compile([]byte("package m\ninteraction S {\n    "+source+"\n}"), base)
		if err != nil || asm != "" || len(ds) != 1 || ds[0].Line != 3 {
			t.Fatalf("bad managed condition accepted: %s err=%v ds=%v", source, err, ds)
		}
	}
}

func TestMessagePredicatesAndDialogueHavePreciseDiagnostics(t *testing.T) {
	for _, tc := range []struct {
		source, message string
		col             int
	}{
		{"dialogue(256)", "numeric value", 14},
		{"dialogue(\"hello\")", "message id", 14},
		{"dialogue(VAR_X)", "message id", 14},
		{"if ask(256) {}", "predicate ask arg 1 numeric value", 12},
		{"if ask(\"hello\") {}", "predicate ask arg 1 requires message id", 12},
		{"if ask() {}", "predicate ask expects 1 args", 8},
		{"if ask(1, 2) {}", "predicate ask expects 1 args", 8},
		{"ask(1)", "ask(message) is a predicate", 5},
		{"if x := 1; flag(FLAG_X) {}", "unsupported if initializer", 8},
	} {
		t.Run(tc.source, func(t *testing.T) {
			asm, ds, err := Compile([]byte("package m\nscript S {\n    "+tc.source+"\n}"), hgssT(t))
			if err != nil || asm != "" || len(ds) != 1 || ds[0].Line != 3 || ds[0].Col != tc.col || !strings.Contains(ds[0].Msg, tc.message) {
				t.Fatalf("want 3:%d %s: err=%v asm=%q diags=%v", tc.col, tc.message, err, asm, ds)
			}
		})
	}
}

// This small command interpreter checks observable branch and lifecycle effects,
// independently of assembly spelling and generated label names. Inputs are the
// HGSS menu's result: 0 = Yes, 1 = No or cancellation.
func playInteraction(t *testing.T, instructions []Instr, tgt target.GameTarget, responses []int, initialFlags map[string]bool) ([]string, map[string]bool) {
	t.Helper()
	labels := map[string]int{}
	for i, in := range instructions {
		if in.Label != "" {
			labels[in.Label] = i
		}
	}
	flags := map[string]bool{}
	for name, value := range initialFlags {
		flags[name] = value
	}
	vars := map[string]int{}
	var effects []string
	locked, messageOpen, waiting := false, false, false
	equal := false
	response := 0
	for pc, steps := 0, 0; pc < len(instructions); pc, steps = pc+1, steps+1 {
		if steps > 100 {
			t.Fatal("flow failed to terminate")
		}
		in := instructions[pc]
		if in.Label != "" || in.Directive != "" {
			continue
		}
		name := in.Macro
		if name == "" {
			name = tgt.MacroName(in.Opcode)
		}
		switch name {
		case "lockall":
			if locked {
				t.Fatal("double lock")
			}
			locked = true
		case "faceplayer":
			if !locked {
				t.Fatal("face without lock")
			}
		case "npc_msg":
			if !locked || messageOpen {
				t.Fatal("invalid message lifecycle")
			}
			messageOpen = true
			effects = append(effects, in.Args[0])
		case "WaitButton":
			if !messageOpen {
				t.Fatal("wait without message")
			}
			waiting = true
		case "yesno":
			if !messageOpen || response >= len(responses) {
				t.Fatal("menu without prompt or test response")
			}
			vars[in.Args[0]] = responses[response]
			response++
			waiting = true
		case "closemsg":
			if !messageOpen || !waiting {
				t.Fatal("close before input")
			}
			messageOpen, waiting = false, false
		case "compare_var_to_value":
			if in.Args[1] != "0" {
				t.Fatalf("unexpected comparator %v", in.Args)
			}
			equal = vars[in.Args[0]] == 0
		case "goto_if_ne":
			if !equal {
				pc = labels[in.Args[0]]
			}
		case "goto_if_unset":
			if !flags[in.Args[0]] {
				pc = labels[in.Args[1]]
			}
		case "goto":
			pc = labels[in.Args[0]]
		case "setflag":
			flags[in.Args[0]] = true
		case "clearflag":
			flags[in.Args[0]] = false
		case "releaseall":
			if !locked || messageOpen {
				t.Fatal("release with leaked message/lock")
			}
			locked = false
		case "end":
			if locked || messageOpen || response != len(responses) {
				t.Fatal("end leaked lock/message or skipped menu")
			}
			return effects, flags
		default:
			t.Fatalf("unexpected instruction %s %v", name, in.Args)
		}
	}
	t.Fatal("flow fell through without End")
	return nil, nil
}

func TestAskBranchesYesNoCancelAndNestedCleanup(t *testing.T) {
	src := []byte(`package m
interaction Offer {
    if ask(MSG_OFFER) {
        flag_set(FLAG_ACCEPTED)
        if ask(MSG_CONFIRM) {
            dialogue(MSG_CONFIRMED)
        } else {
            flag_clear(FLAG_ACCEPTED)
            dialogue(MSG_CANCELLED)
        }
    } else if flag(FLAG_ACCEPTED) {
        dialogue(MSG_REPEAT)
    } else {
        dialogue(MSG_DECLINED)
    }
}`)
	tgt := hgssT(t)
	root, w, err := Parse(src)
	if err != nil {
		t.Fatal(err)
	}
	prog, ds := Resolve(w, root, tgt)
	if len(ds) > 0 {
		t.Fatal(ds)
	}
	for _, tc := range []struct {
		name      string
		responses []int
		messages  []string
		accepted  bool
	}{
		{"yes", []int{0, 0}, []string{"MSG_OFFER", "MSG_CONFIRM", "MSG_CONFIRMED"}, true},
		{"no", []int{1}, []string{"MSG_OFFER", "MSG_DECLINED"}, false},
		{"B cancellation", []int{1}, []string{"MSG_OFFER", "MSG_DECLINED"}, false},
		{"nested no", []int{0, 1}, []string{"MSG_OFFER", "MSG_CONFIRM", "MSG_CANCELLED"}, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			messages, flags := playInteraction(t, Lower(prog), tgt, tc.responses, nil)
			if !reflect.DeepEqual(messages, tc.messages) || flags["FLAG_ACCEPTED"] != tc.accepted {
				t.Fatalf("messages=%v flags=%v", messages, flags)
			}
		})
	}
}

func TestOfferExampleCompletionAndRepeatVisit(t *testing.T) {
	src, err := os.ReadFile("examples/npc-offer.poryz")
	if err != nil {
		t.Fatal(err)
	}
	tgt := hgssT(t)
	root, w, err := Parse(src)
	if err != nil {
		t.Fatal(err)
	}
	prog, ds := Resolve(w, root, tgt)
	if len(ds) != 0 {
		t.Fatal(ds)
	}
	code := Lower(prog)
	messages, flags := playInteraction(t, code, tgt, []int{0}, nil)
	if !reflect.DeepEqual(messages, []string{"0", "1"}) || !flags["FLAG_UNK_042"] {
		t.Fatalf("accepted visit: messages=%v flags=%v", messages, flags)
	}
	messages, flags = playInteraction(t, code, tgt, nil, flags)
	if !reflect.DeepEqual(messages, []string{"3"}) || !flags["FLAG_UNK_042"] {
		t.Fatalf("repeat visit: messages=%v flags=%v", messages, flags)
	}
	for _, response := range []int{1, 1} { // No and the engine-normalized B cancellation.
		messages, flags = playInteraction(t, code, tgt, []int{response}, nil)
		if !reflect.DeepEqual(messages, []string{"0", "2"}) || flags["FLAG_UNK_042"] {
			t.Fatalf("declined visit: messages=%v flags=%v", messages, flags)
		}
	}
}

func TestInteractionFormatterPreservesSemanticsAndIsStable(t *testing.T) {
	src := []byte("package m\n\ninteraction S {\nif ask(MSG_OFFER) {\ndialogue(MSG_YES) // } is a comment\n} else {\ndialogue(MSG_NO)\n}\n}\n")
	formatted, err := Format(src)
	if err != nil {
		t.Fatal(err)
	}
	want := "package m\n\ninteraction S {\n    if ask(MSG_OFFER) {\n        dialogue(MSG_YES) // } is a comment\n    } else {\n        dialogue(MSG_NO)\n    }\n}"
	if string(formatted) != want {
		t.Fatalf("format:\n%s", formatted)
	}
	again, err := Format(formatted)
	if err != nil || string(again) != want {
		t.Fatal("formatter is not stable", err)
	}
	if compileFlow(t, string(src), hgssT(t)) != compileFlow(t, string(formatted), hgssT(t)) {
		t.Fatal("format changed behavior")
	}
}

func TestRawYesNoExampleUsesEnginePolarity(t *testing.T) {
	src, err := os.ReadFile("examples/yes-no.poryz")
	if err != nil {
		t.Fatal(err)
	}
	asm := compileFlow(t, string(src), hgssT(t))
	if !strings.Contains(asm, "\tCompare VAR_SPECIAL_RESULT, 0\n\tGoToIfEq _Yes\n") {
		t.Fatalf("Yes must branch on result 0: %s", asm)
	}
}
