package poryscriptz

import (
	gts "github.com/odvcencio/gotreesitter"
	"github.com/odvcencio/gotreesitter/taproot"
)

//go:generate go run ./cmd/genblob

// Parse parses a poryscriptZ source file from src, returning the root CST
// node, a Walker for traversing it, and any error.  The embedded language
// blob is used when available (and its hash matches the current grammar);
// otherwise go generate must be re-run.
func Parse(src []byte) (*gts.Node, *taproot.Walker, error) {
	blob, err := loadEmbeddedLanguageBlob(PoryscriptZGrammar())
	if err != nil {
		return nil, nil, err
	}
	return taproot.ParseFromBlob("poryscriptz", blob, PoryscriptZGrammar, src)
}

// firstScriptName walks the CST rooted at root and returns the name field of
// the first script_declaration encountered.
func firstScriptName(w *taproot.Walker, root *gts.Node) string {
	var found string
	var walk func(n *gts.Node)
	walk = func(n *gts.Node) {
		if n == nil || found != "" {
			return
		}
		if w.Type(n) == "script_declaration" {
			found = w.Text(w.Field(n, "name"))
			return
		}
		for i := 0; i < n.ChildCount(); i++ {
			walk(n.Child(i))
		}
	}
	walk(root)
	return found
}
