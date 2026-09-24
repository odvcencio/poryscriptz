package vocab

type ArgKind int

const (
	ArgU8    ArgKind = iota // raw width 1
	ArgU16                  // raw width 2
	ArgU32                  // raw width 4
	ArgVar                  // "var"    → a VAR_* reference
	ArgFlag                 // "flag"   → a FLAG_* reference
	ArgLabel                // "script" → a local jump target
	ArgSym                  // any other semantic type (species/item/move/message/…) — emitted as-is
)

type Command struct {
	Opcode int
	Name   string // canonical name, e.g. "setflag"; "scrcmd_NNN" if unnamed
	Args   []ArgKind
}

type Table struct {
	byName   map[string]*Command
	byOpcode map[int]*Command
}

func (t *Table) ByName(n string) (*Command, bool) { c, ok := t.byName[n]; return c, ok }
func (t *Table) ByOpcode(op int) (*Command, bool) { c, ok := t.byOpcode[op]; return c, ok }
func (t *Table) Len() int                         { return len(t.byOpcode) }
func (t *Table) Entries() []*Command {
	if t == nil {
		return nil
	}
	out := make([]*Command, 0, len(t.byName))
	for _, e := range t.byName {
		out = append(out, e)
	}
	return out
}

// MacroEntry describes a decomp script.inc macro: a name and its argument
// count/shapes. Macro calls in .poryz emit verbatim into the .s output; the
// decomp assembler expands them via script.inc at assemble time.
type MacroEntry struct {
	Name      string
	EmitName  string
	Args      []ArgKind
	MinArgs   int
	Doc       string
	Params    []string
	Movement  bool
	Canonical string
	Opcode    *int
}

// MacroTable is the lookup table for decomp macros (loaded from macros.json).
type MacroTable struct {
	byName map[string]*MacroEntry
}

func (m *MacroTable) ByName(n string) (*MacroEntry, bool) {
	if m == nil {
		return nil, false
	}
	e, ok := m.byName[n]
	return e, ok
}

func (m *MacroTable) Len() int {
	if m == nil {
		return 0
	}
	return len(m.byName)
}

func (m *MacroTable) Entries() []*MacroEntry {
	if m == nil {
		return nil
	}
	out := make([]*MacroEntry, 0, len(m.byName))
	for _, e := range m.byName {
		out = append(out, e)
	}
	return out
}

func (m *MacroTable) Add(e *MacroEntry) {
	if m.byName == nil {
		m.byName = map[string]*MacroEntry{}
	}
	m.byName[e.Name] = e
}
