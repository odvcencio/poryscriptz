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
