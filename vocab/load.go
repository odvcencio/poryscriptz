package vocab

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
)

//go:embed testdata/*.json
var bundled embed.FS

func Bundled(name string) ([]byte, error) { return bundled.ReadFile("testdata/" + name) }

type rawFile struct {
	Commands []rawCmd `json:"commands"`
}
type rawCmd struct {
	Name string            `json:"name"`
	Args []json.RawMessage `json:"args"`
}

type rawMacroFile struct {
	Macros []rawMacroEntry `json:"macros"`
}
type rawMacroEntry struct {
	Name      string   `json:"name"`
	Args      []string `json:"args"`
	MinArgs   *int     `json:"min_args,omitempty"`
	Doc       string   `json:"doc,omitempty"`
	Params    []string `json:"params,omitempty"`
	Movement  bool     `json:"movement,omitempty"`
	Canonical string   `json:"canonical,omitempty"`
	Opcode    *int     `json:"opcode,omitempty"`
}

func LoadScrcmdJSON(path string) (*Table, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read scrcmd.json: %w", err)
	}
	return ParseScrcmdJSON(b)
}

func ParseScrcmdJSON(b []byte) (*Table, error) {
	var rf rawFile
	if err := json.Unmarshal(b, &rf); err != nil {
		return nil, fmt.Errorf("parse scrcmd.json: %w", err)
	}
	t := &Table{byName: map[string]*Command{}, byOpcode: map[int]*Command{}}
	for op, rc := range rf.Commands {
		c := &Command{Opcode: op, Name: rc.Name, Args: make([]ArgKind, 0, len(rc.Args))}
		for i, raw := range rc.Args {
			ak, err := argKind(raw)
			if err != nil {
				return nil, fmt.Errorf("cmd %q arg %d: %w", rc.Name, i, err)
			}
			c.Args = append(c.Args, ak)
		}
		t.byName[c.Name] = c
		t.byOpcode[op] = c
	}
	return t, nil
}

// LoadMacrosJSON reads a macros.json file and returns a MacroTable.
func LoadMacrosJSON(path string) (*MacroTable, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read macros.json: %w", err)
	}
	return ParseMacrosJSON(b)
}

func ParseMacrosJSON(b []byte) (*MacroTable, error) {
	var rf rawMacroFile
	if err := json.Unmarshal(b, &rf); err != nil {
		return nil, fmt.Errorf("parse macros.json: %w", err)
	}
	mt := &MacroTable{byName: map[string]*MacroEntry{}}
	for _, rm := range rf.Macros {
		e := &MacroEntry{Name: rm.Name, Args: make([]ArgKind, 0, len(rm.Args)), MinArgs: len(rm.Args), Doc: rm.Doc, Params: rm.Params, Movement: rm.Movement, Canonical: rm.Canonical, Opcode: rm.Opcode}
		if rm.MinArgs != nil {
			e.MinArgs = *rm.MinArgs
		}
		for i, s := range rm.Args {
			ak, err := macroArgKind(s)
			if err != nil {
				return nil, fmt.Errorf("macro %q arg %d: %w", rm.Name, i, err)
			}
			e.Args = append(e.Args, ak)
		}
		mt.byName[e.Name] = e
	}
	return mt, nil
}

func macroArgKind(s string) (ArgKind, error) {
	switch s {
	case "flag": // reserved for future macros that take FLAG_* arguments; no current macro uses this
		return ArgFlag, nil
	case "var":
		return ArgVar, nil
	case "script": // reserved for future macros that take label/script arguments; no current macro uses this
		return ArgLabel, nil
	case "sym", "":
		return ArgSym, nil
	default:
		return ArgSym, fmt.Errorf("unrecognised macro arg kind %q", s)
	}
}

func argKind(raw json.RawMessage) (ArgKind, error) {
	var n int
	if json.Unmarshal(raw, &n) == nil {
		switch n {
		case 1:
			return ArgU8, nil
		case 4:
			return ArgU32, nil
		default:
			return ArgU16, nil
		}
	}
	var s string
	if json.Unmarshal(raw, &s) == nil {
		switch s {
		case "flag":
			return ArgFlag, nil
		case "var":
			return ArgVar, nil
		case "script":
			return ArgLabel, nil
		default:
			return ArgSym, nil
		}
	}
	return ArgSym, fmt.Errorf("unrecognised arg token %s", string(raw))
}
