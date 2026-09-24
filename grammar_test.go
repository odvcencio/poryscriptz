package poryscriptz

import "testing"

func TestParseScriptDecl(t *testing.T) {
	src := []byte("package m\n\nscript YoungsterJoey {\n\tsetflag(FLAG_X)\n}\n")
	root, w, err := Parse(src)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if root.HasError() {
		t.Fatalf("unexpected parse error: %v", w.SyntaxError(root))
	}
	if got := firstScriptName(w, root); got != "YoungsterJoey" {
		t.Fatalf("script name: want YoungsterJoey, got %q", got)
	}
}
