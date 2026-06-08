package vocab

type ArgKind int

const (
	ArgU8  ArgKind = iota // raw width 1
	ArgU16                // raw width 2
	ArgU32                // raw width 4
	ArgVar                // "var"    → a VAR_* reference
	ArgFlag               // "flag"   → a FLAG_* reference
	ArgLabel              // "script" → a local jump target
	ArgSym                // any other semantic type (species/item/move/message/…) — emitted as-is
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
func (t *Table) ByOpcode(op int) (*Command, bool)  { c, ok := t.byOpcode[op]; return c, ok }
func (t *Table) Len() int                          { return len(t.byOpcode) }
