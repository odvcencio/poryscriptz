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
