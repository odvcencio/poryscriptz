package poryscriptz_test

import (
	"strings"
	"testing"

	poryz "github.com/odvcencio/poryscriptz"
	"github.com/odvcencio/poryscriptz/target"
)

func TestDecompileRoundTripLayouts(t *testing.T) {
	tgt, err := target.HGSS("")
	if err != nil {
		t.Fatal(err)
	}
	cases := []string{
		"#include \"constants/scrcmd.h\"\n\t.include \"asm/macros/script.inc\"\n\n\t.rodata\n\tScrDef A\n\tScrDefEnd\nA:\n\tEnd\n\t.balign 4, 0\n",
		"#include \"constants/scrcmd.h\"\n\t.include \"asm/macros/script.inc\"\n\n\t.rodata\n\t.option alignment off\n\tInitScriptEntry_OnResume 1\n\tInitScriptEntryEnd\n\tInitScriptEnd\n",
	}
	for _, src := range cases {
		code, err := poryz.Decompile([]byte(src), tgt)
		if err != nil {
			t.Fatal(err)
		}
		asm, diags, err := poryz.Compile([]byte(code), tgt)
		if err != nil {
			t.Fatal(err)
		}
		if len(diags) > 0 {
			t.Fatal(diags)
		}
		if !strings.Contains(asm, "\tEnd") && !strings.Contains(asm, "\tInitScriptEnd") {
			t.Fatalf("lost commands: %s", asm)
		}
	}
}
func TestDecompileUnknownMacroHasFixHint(t *testing.T) {
	tgt, err := target.HGSS("")
	if err != nil {
		t.Fatal(err)
	}
	_, err = poryz.Decompile([]byte("#include \"constants/scrcmd.h\"\n\t.include \"asm/macros/script.inc\"\n\t.rodata\n\tScrDefEnd\nA:\n\tUnknownThing\n"), tgt)
	if err == nil || !strings.Contains(err.Error(), "6:2") || !strings.Contains(err.Error(), "regenerate the vocabulary") {
		t.Fatalf("want positioned fix hint, got %v", err)
	}
}

func TestDecompileRejectsWrongKnownArity(t *testing.T) {
	tgt, err := target.HGSS("")
	if err != nil {
		t.Fatal(err)
	}
	_, err = poryz.Decompile([]byte("#include \"constants/scrcmd.h\"\n\t.include \"asm/macros/script.inc\"\n\t.rodata\n\tScrDefEnd\nA:\n\tEnd extra\n"), tgt)
	if err == nil || !strings.Contains(err.Error(), "6:2") || !strings.Contains(err.Error(), "expects 0 to 0 args") {
		t.Fatalf("want positioned arity error, got %v", err)
	}
}
func TestFormatIdempotent(t *testing.T) {
	src := []byte("package m\nscript A {\nif flag(FLAG_UNK_042) {\nEnd()\n}\n}\n")
	once, err := poryz.Format(src)
	if err != nil {
		t.Fatal(err)
	}
	twice, err := poryz.Format(once)
	if err != nil {
		t.Fatal(err)
	}
	if string(once) != string(twice) {
		t.Fatalf("format changed twice: %q / %q", once, twice)
	}
}
