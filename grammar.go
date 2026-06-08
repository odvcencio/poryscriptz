package poryscriptz

import (
	gg "github.com/odvcencio/gotreesitter/grammargen"
)

// Type aliases so language_blob_embed.go and cmd/genblob compile without
// importing grammargen directly.
type Grammar = gg.Grammar
type Rule = gg.Rule

// Re-exports used by cmd/genblob/main.go.
var ExportGrammarJSON = gg.ExportGrammarJSON
var GenerateLanguageAndBlob = gg.GenerateLanguageAndBlob

func PoryscriptZGrammar() *Grammar {
	return gg.ExtendGrammar("poryscriptz", gg.GoGrammar(), func(g *Grammar) {
		g.Define("script_declaration",
			gg.Seq(
				gg.Str("script"),
				gg.Field("name", gg.Sym("identifier")),
				gg.Sym("block"),
			))
		gg.AppendChoice(g, "_top_level_declaration", gg.Sym("script_declaration"))
		gg.AddConflict(g, "_top_level_declaration", "script_declaration")
	})
}
