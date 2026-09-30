package poryscriptz

import (
	"strings"
	"testing"
)

func TestUnsupportedStatementsNeverDisappear(t *testing.T) {
	for _, statement := range []string{
		"break", "continue", "return", "goto _Done", "goto(_Done)",
		"VAR_X = 1", "x := 1", "1 + 2", "{ End() }", "defer End()",
	} {
		t.Run(statement, func(t *testing.T) {
			asm, diags, err := Compile([]byte("package m\nscript S {\n    "+statement+"\n}\nlabel _Done {}\n"), hgssT(t))
			if err != nil {
				t.Fatal(err)
			}
			if asm != "" || len(diags) != 1 || diags[0].Line != 3 || diags[0].Col != 5 || !strings.Contains(diags[0].Msg, "unsupported statement") {
				t.Fatalf("unsupported source must fail at 3:5, without output: asm=%q diags=%v", asm, diags)
			}
		})
	}
}

func TestUnsupportedTopLevelDeclarationsNeverDisappear(t *testing.T) {
	for _, declaration := range []string{
		"func unused() { End() }", "var ignored int", "const ignored = 1", "type Ignored int", "import \"fmt\"",
	} {
		t.Run(declaration, func(t *testing.T) {
			asm, diags, err := Compile([]byte("package m\n"+declaration+"\nscript S { End() }\n"), hgssT(t))
			if err != nil {
				t.Fatal(err)
			}
			if asm != "" || len(diags) != 1 || diags[0].Line != 2 || diags[0].Col != 1 || !strings.Contains(diags[0].Msg, "unsupported top-level") {
				t.Fatalf("unsupported declaration must fail at 2:1: asm=%q diags=%v", asm, diags)
			}
		})
	}
}

func TestRawControlFlowEscapeHatchesRemainAvailable(t *testing.T) {
	src := []byte("package m\nscript S { AsmGoto(_Done) }\nlabel _Done { AsmReturn() }\n")
	asm, diags, err := Compile(src, hgssT(t))
	if err != nil || len(diags) != 0 || !strings.Contains(asm, "\tgoto _Done\n") || !strings.Contains(asm, "\treturn\n") {
		t.Fatalf("raw escapes changed: err=%v diags=%v asm=%s", err, diags, asm)
	}
}
