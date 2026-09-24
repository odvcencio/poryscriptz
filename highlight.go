package poryscriptz

import (
	_ "embed"
	"sync"

	gts "github.com/odvcencio/gotreesitter"
	"github.com/odvcencio/gotreesitter/taproot"
)

//go:embed queries/highlights.scm
var highlightQuery string

var highlightOnce sync.Once
var highlighter *gts.Highlighter
var highlightErr error

// Highlight returns grammar-derived syntax spans for editors.
func Highlight(src []byte) ([]gts.HighlightRange, error) {
	highlightOnce.Do(func() {
		blob, err := loadEmbeddedLanguageBlob(PoryscriptZGrammar())
		if err != nil {
			highlightErr = err
			return
		}
		lang, err := taproot.LanguageFromBlob("poryscriptz", blob, PoryscriptZGrammar)
		if err != nil {
			highlightErr = err
			return
		}
		highlighter, highlightErr = gts.NewHighlighter(lang, highlightQuery)
	})
	if highlightErr != nil {
		return nil, highlightErr
	}
	return highlighter.Highlight(src), nil
}
