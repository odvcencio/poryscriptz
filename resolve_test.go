package poryscriptz

import (
	"strings"
	"testing"

	"m31labs.dev/poryscriptz/target"
)

func hgssT(t *testing.T) target.GameTarget {
	t.Helper()
	g, err := target.HGSS("vocab/testdata/scrcmd.json")
	if err != nil {
		t.Fatalf("HGSS: %v", err)
	}
	return g
}

func TestResolveUnknownCommand(t *testing.T) {
	src := []byte("package m\n\nscript S {\n\tdefinitely_not_a_command(1)\n}\n")
	_, diags, _ := Compile(src, hgssT(t))
	if len(diags) == 0 || !strings.Contains(diags[0].Msg, "not in") {
		t.Fatalf("want unknown-command diagnostic, got %v", diags)
	}
}

func TestResolveKnownCommandOK(t *testing.T) {
	src := []byte("package m\n\nscript S {\n\tsetflag(FLAG_X)\n\tend()\n}\n")
	_, diags, err := Compile(src, hgssT(t))
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	if len(diags) != 0 {
		t.Fatalf("want no diags, got %v", diags)
	}
}

func TestResolveArityMismatch(t *testing.T) {
	src := []byte("package m\n\nscript S {\n\tsetflag()\n}\n")
	_, diags, err := Compile(src, hgssT(t))
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	if len(diags) == 0 || !strings.Contains(diags[0].Msg, "expects") {
		t.Fatalf("want arity-mismatch diagnostic, got %v", diags)
	}
}

func TestResolveElseAccepted(t *testing.T) {
	src := []byte("package m\nscript S {\n\tif flag(FLAG_X) {\n\t\tend()\n\t} else {\n\t\tend()\n\t}\n}\n")
	_, diags, err := Compile(src, hgssT(t))
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	if len(diags) != 0 {
		t.Fatalf("want no diags for else clause, got %v", diags)
	}
}

func TestResolveMacroCallOK(t *testing.T) {
	src := []byte("package m\nscript S {\n\tnpc_msg(5)\n\tclosemsg()\n\tend()\n}\n")
	_, diags, err := Compile(src, hgssT(t))
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	if len(diags) != 0 {
		t.Fatalf("want no diags for macro calls, got %v", diags)
	}
}

func TestResolveMacroArityMismatch(t *testing.T) {
	// trainer_battle expects 4 args; supply 2.
	src := []byte("package m\nscript S {\n\ttrainer_battle(900, 0)\n}\n")
	_, diags, err := Compile(src, hgssT(t))
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	if len(diags) == 0 || !strings.Contains(diags[0].Msg, "expects") {
		t.Fatalf("want arity-mismatch diagnostic, got %v", diags)
	}
}

func TestResolveUnknownStillErrors(t *testing.T) {
	// Typo: "npc_msg_typo" — must not silently pass.
	src := []byte("package m\nscript S {\n\tnpc_msg_typo(5)\n}\n")
	_, diags, _ := Compile(src, hgssT(t))
	if len(diags) == 0 || !strings.Contains(diags[0].Msg, "not in") {
		t.Fatalf("want unknown-command diagnostic, got %v", diags)
	}
}

// TestResolveMacroCallViaExclusiveName proves that giveitem_no_check — the only
// macro name that does NOT collide with any scrcmd command — flows through the
// macro-table branch of resolveCall and produces a MacroCall stmt, not a Call.
func TestResolveMacroCallViaExclusiveName(t *testing.T) {
	src := []byte("package m\nscript S {\n\tgiveitem_no_check(ITEM_POTION, 1)\n\tend()\n}\n")
	root, w, err := Parse(src)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	prog, diags := Resolve(w, root, hgssT(t))
	if len(diags) != 0 {
		t.Fatalf("want no diags, got %v", diags)
	}
	if len(prog.Scripts) == 0 || len(prog.Scripts[0].Stmts) == 0 {
		t.Fatalf("expected at least one stmt, got prog=%v", prog)
	}
	// The first stmt must be a MacroCall, not a Call.
	mc, ok := prog.Scripts[0].Stmts[0].(MacroCall)
	if !ok {
		t.Fatalf("want MacroCall for giveitem_no_check, got %T", prog.Scripts[0].Stmts[0])
	}
	if mc.Macro.Name != "giveitem_no_check" {
		t.Fatalf("want macro name giveitem_no_check, got %q", mc.Macro.Name)
	}
	if len(mc.Args) != 2 {
		t.Fatalf("want 2 args, got %d", len(mc.Args))
	}
}

// TestResolveMacroArityMismatchExclusive proves that calling giveitem_no_check
// with the wrong arity produces the MACRO-branch diagnostic (prefix "macro …
// expects"), not the scrcmd-branch message.  This is distinct from
// TestResolveMacroArityMismatch which covers a scrcmd-colliding name.
func TestResolveMacroArityMismatchExclusive(t *testing.T) {
	// giveitem_no_check expects 2 args; supply 1.
	src := []byte("package m\nscript S {\n\tgiveitem_no_check(ITEM_POTION)\n}\n")
	_, diags, err := Compile(src, hgssT(t))
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	if len(diags) == 0 || !strings.Contains(diags[0].Msg, "macro giveitem_no_check expects") {
		t.Fatalf("want macro-branch arity diagnostic, got %v", diags)
	}
}

// TestResolveScrCmdWinsOnCollision pins the deliberate "scrcmd wins on
// collision" order: give_mon appears in both scrcmd.json (6 args) and
// macros.json (6 args).  Supplying 6 correct args must resolve to a scrcmd
// Call, never a MacroCall, so that future reordering cannot silently change
// the emitted bytecode.
func TestResolveScrCmdWinsOnCollision(t *testing.T) {
	src := []byte("package m\nscript S {\n\tgive_mon(SPECIES_PIKACHU, 5, 0, 0, 0, VAR_RESULT)\n\tend()\n}\n")
	root, w, err := Parse(src)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	prog, diags := Resolve(w, root, hgssT(t))
	if len(diags) != 0 {
		t.Fatalf("want no diags, got %v", diags)
	}
	if len(prog.Scripts) == 0 || len(prog.Scripts[0].Stmts) == 0 {
		t.Fatalf("expected at least one stmt")
	}
	// Must be a Call (scrcmd), NOT a MacroCall.
	if _, ok := prog.Scripts[0].Stmts[0].(Call); !ok {
		t.Fatalf("want scrcmd Call for give_mon (scrcmd wins on collision), got %T", prog.Scripts[0].Stmts[0])
	}
	if _, ok := prog.Scripts[0].Stmts[0].(MacroCall); ok {
		t.Fatalf("give_mon must NOT resolve as MacroCall — scrcmd must win on collision")
	}
}
