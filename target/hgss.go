package target

import (
	"fmt"
	"path/filepath"
	"strings"

	"m31labs.dev/poryscriptz/vocab"
)

type hgss struct {
	tbl    *vocab.Table
	macros *vocab.MacroTable
}

// HGSS loads the HGSS game target from scrcmdJSONPath.
// It also attempts to load a macros.json file from the same directory as
// scrcmdJSONPath; if that file does not exist, the macro table is empty
// (all macros will be reported as unknown identifiers).
func HGSS(scrcmdJSONPath string) (GameTarget, error) {
	tbl, err := vocab.LoadScrcmdJSON(scrcmdJSONPath)
	if err != nil {
		return nil, err
	}
	macrosPath := filepath.Join(filepath.Dir(scrcmdJSONPath), "macros.json")
	mt, err := vocab.LoadMacrosJSON(macrosPath)
	if err != nil {
		// macros.json is optional; treat a missing file as an empty table.
		mt = nil
	}
	return &hgss{tbl: tbl, macros: mt}, nil
}

func (h *hgss) Name() string              { return "HGSS" }
func (h *hgss) Vocabulary() *vocab.Table  { return h.tbl }
func (h *hgss) Macros() *vocab.MacroTable { return h.macros }

func (h *hgss) MacroName(opcode int) string {
	if c, ok := h.tbl.ByOpcode(opcode); ok && c.Name != "" {
		return c.Name
	}
	return fmt.Sprintf("scrcmd_%d", opcode)
}

// Preamble matches the header of every real scr_seq/*.s (lines at column 0).
func (h *hgss) Preamble(headers ...string) string {
	var b strings.Builder
	b.WriteString("#include \"constants/scrcmd.h\"\n")
	for _, header := range headers {
		fmt.Fprintf(&b, "#include \"%s\"\n", header)
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
