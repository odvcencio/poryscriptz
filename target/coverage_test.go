package target

import "testing"

func TestBundledVocabularyCoverage(t *testing.T) {
	game, err := HGSS("")
	if err != nil {
		t.Fatal(err)
	}
	h := game.(*hgss)
	if h.tbl.Len() != 853 {
		t.Fatalf("want 853 opcodes, got %d", h.tbl.Len())
	}
	canonical := 0
	for _, e := range h.macros.Entries() {
		if e.Canonical == "" && e.Name != "AsmGoto" && e.Name != "AsmReturn" && e.Name != "AsmSwitch" && e.Name != "AsmCase" {
			canonical++
		}
	}
	if canonical < 1015 {
		t.Fatalf("want 1015 public macros, got %d", canonical)
	}
	for op := 0; op < h.tbl.Len(); op++ {
		if h.CanonicalOpcode(op) == "" {
			t.Fatalf("opcode %d has no macro", op)
		}
	}
	if got := h.CanonicalOpcode(0); got != "Noop" {
		t.Fatalf("opcode 0 mapped to %q", got)
	}
}
