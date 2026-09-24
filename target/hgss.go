package target

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/odvcencio/poryscriptz/vocab"
)

type hgss struct {
	tbl               *vocab.Table
	macros            *vocab.MacroTable
	defaultPret       bool
	canonicalByOpcode map[int]string
	canonicalByName   map[string]string
}

// HGSS loads the HGSS game target from scrcmdJSONPath.
// It also attempts to load a macros.json file from the same directory as
// scrcmdJSONPath; if that file does not exist, the macro table is empty
// (all macros will be reported as unknown identifiers).
func HGSS(scrcmdJSONPath string) (GameTarget, error) {
	var tbl *vocab.Table
	var err error
	if scrcmdJSONPath == "" {
		var data []byte
		data, err = vocab.Bundled("scrcmd.json")
		if err == nil {
			tbl, err = vocab.ParseScrcmdJSON(data)
		}
	} else {
		tbl, err = vocab.LoadScrcmdJSON(scrcmdJSONPath)
	}
	if err != nil {
		return nil, err
	}
	var mt *vocab.MacroTable
	if scrcmdJSONPath == "" {
		data, e := vocab.Bundled("macros.json")
		if e == nil {
			mt, err = vocab.ParseMacrosJSON(data)
		} else {
			err = e
		}
	} else {
		mt, err = vocab.LoadMacrosJSON(filepath.Join(filepath.Dir(scrcmdJSONPath), "macros.json"))
	}
	if err != nil {
		// macros.json is optional; treat a missing file as an empty table.
		mt = nil
	}
	var generated *vocab.MacroTable
	if scrcmdJSONPath == "" {
		data, e := vocab.Bundled("decomp_macros.json")
		if e == nil {
			generated, err = vocab.ParseMacrosJSON(data)
		} else {
			err = e
		}
	} else {
		generated, err = vocab.LoadMacrosJSON(filepath.Join(filepath.Dir(scrcmdJSONPath), "decomp_macros.json"))
	}
	if err == nil {
		if mt == nil {
			mt = generated
		} else {
			for _, e := range generated.Entries() {
				mt.Add(e)
			}
		}
	}
	if mt != nil {
		for alias, raw := range map[string]string{"AsmGoto": "goto", "AsmReturn": "return", "AsmSwitch": "switch", "AsmCase": "case"} {
			if entry, ok := mt.ByName(raw); ok {
				copy := *entry
				copy.Name, copy.EmitName = alias, raw
				mt.Add(&copy)
			}
		}
	}
	h := &hgss{tbl: tbl, macros: mt, defaultPret: scrcmdJSONPath == "", canonicalByOpcode: map[int]string{}, canonicalByName: map[string]string{}}
	entries := mt.Entries()
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name < entries[j].Name })
	for _, e := range entries {
		if e.Canonical != "" {
			continue
		}
		normalized := strings.ToLower(strings.ReplaceAll(e.Name, "_", ""))
		if _, exists := h.canonicalByName[normalized]; !exists {
			h.canonicalByName[normalized] = e.Name
		}
		if e.Opcode != nil && !strings.HasPrefix(e.Name, "InitScript") {
			if _, exists := h.canonicalByOpcode[*e.Opcode]; !exists {
				h.canonicalByOpcode[*e.Opcode] = e.Name
			}
		}
	}
	return h, nil
}

func (h *hgss) DefaultPret() bool { return h.defaultPret }
func (h *hgss) CanonicalMacro(name string) string {
	if entry, ok := h.macros.ByName(name); ok && entry.Canonical != "" {
		return entry.Canonical
	}
	if len(name) > 0 && name[0] >= 'A' && name[0] <= 'Z' {
		return name
	}
	if canonical, ok := h.canonicalByName[strings.ToLower(strings.ReplaceAll(name, "_", ""))]; ok {
		return canonical
	}
	return name
}
func (h *hgss) CanonicalOpcode(opcode int) string {
	if name, ok := h.canonicalByOpcode[opcode]; ok {
		return name
	}
	return h.CanonicalMacro(h.MacroName(opcode))
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
