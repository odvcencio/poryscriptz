package target

import "testing"

func TestHGSSVocabularyAndMacros(t *testing.T) {
	hg, err := HGSS("../vocab/testdata/scrcmd.json")
	if err != nil {
		t.Fatalf("HGSS: %v", err)
	}
	if _, ok := hg.Vocabulary().ByName("setflag"); !ok {
		t.Fatal("setflag missing from HGSS vocabulary")
	}
	// Unnamed opcodes fall back to scrcmd_NNN.
	if got := hg.MacroName(9999); got != "scrcmd_9999" {
		t.Fatalf("fallback macro: want scrcmd_9999, got %q", got)
	}
}
