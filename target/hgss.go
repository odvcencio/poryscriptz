package target

import (
	"fmt"
	"strings"

	"m31labs.dev/poryscriptz/vocab"
)

type hgss struct{ tbl *vocab.Table }

func HGSS(scrcmdJSONPath string) (GameTarget, error) {
	tbl, err := vocab.LoadScrcmdJSON(scrcmdJSONPath)
	if err != nil {
		return nil, err
	}
	return &hgss{tbl: tbl}, nil
}

func (h *hgss) Vocabulary() *vocab.Table { return h.tbl }

func (h *hgss) MacroName(opcode int) string {
	if c, ok := h.tbl.ByOpcode(opcode); ok && c.Name != "" {
		return c.Name
	}
	return fmt.Sprintf("scrcmd_%d", opcode)
}

// Preamble matches the header of every real scr_seq/*.s (lines at column 0).
func (h *hgss) Preamble(eventHeader string) string {
	var b strings.Builder
	b.WriteString("#include \"constants/scrcmd.h\"\n")
	if eventHeader != "" {
		fmt.Fprintf(&b, "#include \"%s\"\n", eventHeader)
	}
	b.WriteString(".include \"asm/macros/script.inc\"\n\n.rodata\n")
	return b.String()
}

func (h *hgss) Header(labels []string) string {
	var b strings.Builder
	for _, l := range labels {
		fmt.Fprintf(&b, "\tscrdef %s\n", l)
	}
	b.WriteString("\tscrdef_end\n")
	return b.String()
}
