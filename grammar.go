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
		// `script` is a hard keyword ONLY at top-level-declaration scope; the GLR
		// conflict registered below lets it fall through to `identifier` in statements
		// and expressions. (Latent collision only for a top-level `func script()`,
		// which poryscriptZ never emits.)
		g.Define("script_declaration",
			gg.Seq(
				gg.Str("script"),
				gg.Field("name", gg.Sym("identifier")),
				gg.Sym("block"),
			))
		g.Define("label_declaration", gg.Seq(
			gg.Str("label"),
			gg.Field("name", gg.Sym("identifier")),
			gg.Sym("block"),
		))
		g.Define("movement_declaration", gg.Seq(
			gg.Str("movement"),
			gg.Field("name", gg.Sym("identifier")),
			gg.Sym("block"),
		))
		g.Define("include_declaration", gg.Seq(
			gg.Str("include"),
			gg.Field("path", gg.Sym("interpreted_string_literal")),
		))
		g.Define("align_declaration", gg.Seq(
			gg.Str("align4"), gg.Str("("), gg.Str(")"),
		))
		for _, rule := range []string{
			"script_declaration", "label_declaration", "movement_declaration", "include_declaration", "align_declaration",
		} {
			gg.AppendChoice(g, "_top_level_declaration", gg.Sym(rule))
			gg.AddConflict(g, "_top_level_declaration", rule)
		}
	})
}
