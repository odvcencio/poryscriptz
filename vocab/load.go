package vocab

import (
	"encoding/json"
	"fmt"
	"os"
)

type rawFile struct {
	Commands []rawCmd `json:"commands"`
}
type rawCmd struct {
	Name string            `json:"name"`
	Args []json.RawMessage `json:"args"`
}

func LoadScrcmdJSON(path string) (*Table, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read scrcmd.json: %w", err)
	}
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
