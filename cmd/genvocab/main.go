// genvocab derives the command reference from a pokeheartgold checkout.
// Usage: go run ./cmd/genvocab <decomp-root>
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

type entry struct {
	Name      string   `json:"name"`
	Args      []string `json:"args"`
	Params    []string `json:"params,omitempty"`
	MinArgs   int      `json:"min_args"`
	Doc       string   `json:"doc,omitempty"`
	Movement  bool     `json:"movement,omitempty"`
	Canonical string   `json:"canonical,omitempty"`
	Opcode    *int     `json:"opcode,omitempty"`
}

func collect(path string, movement bool) ([]entry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var out []entry
	current := -1
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == ".endm" {
			current = -1
			continue
		}
		if current >= 0 && !movement && out[current].Opcode == nil && strings.HasPrefix(line, ".short ") {
			fields := strings.Fields(strings.TrimPrefix(line, ".short "))
			if len(fields) > 0 {
				if n, err := strconv.Atoi(fields[0]); err == nil && n >= 0 && n < 10000 {
					out[current].Opcode = &n
				}
			}
		}
		if strings.HasPrefix(line, ";") {
			continue
		}
		if !strings.HasPrefix(line, ".macro ") {
			continue
		}
		line = strings.SplitN(strings.TrimSpace(strings.TrimPrefix(line, ".macro ")), ";", 2)[0]
		name, rest, _ := strings.Cut(strings.TrimSpace(line), " ")
		e := entry{Name: name, Args: []string{}, Params: []string{}, Movement: movement}
		if rest != "" {
			for _, arg := range strings.Split(rest, ",") {
				arg = strings.TrimSpace(arg)
				if arg == "" {
					continue
				}
				if !strings.Contains(arg, "=") {
					e.MinArgs++
				}
				e.Args = append(e.Args, "sym")
				e.Params = append(e.Params, arg)
			}
		}
		out = append(out, e)
		current = len(out) - 1
	}
	return out, s.Err()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: genvocab <decomp-root>")
		os.Exit(2)
	}
	root := os.Args[1]
	script, err := collect(filepath.Join(root, "asm/macros/script.inc"), false)
	if err != nil {
		panic(err)
	}
	movement, err := collect(filepath.Join(root, "asm/macros/movement.inc"), true)
	if err != nil {
		panic(err)
	}
	all := append(script, movement...)
	// A few authored files pass a trailing operand that the macro declaration
	// does not name. Keep those calls expressible by measuring real usage.
	usage := map[string]int{}
	files, _ := filepath.Glob(filepath.Join(root, "files/fielddata/script/scr_seq/*.s"))
	for _, path := range files {
		f, err := os.Open(path)
		if err != nil {
			panic(err)
		}
		s := bufio.NewScanner(f)
		for s.Scan() {
			line := strings.TrimSpace(strings.SplitN(s.Text(), ";", 2)[0])
			name, rest, ok := strings.Cut(line, " ")
			if !ok || strings.HasPrefix(name, ".") {
				continue
			}
			n := 1
			for _, c := range rest {
				if c == ',' {
					n++
				}
			}
			if n > usage[name] {
				usage[name] = n
			}
		}
		f.Close()
	}
	for i := range all {
		for len(all[i].Args) < usage[all[i].Name] {
			all[i].Args = append(all[i].Args, "sym")
		}
	}
	for _, e := range append([]entry(nil), all...) {
		alias := snake(e.Name)
		if alias != e.Name {
			e.Canonical = e.Name
			e.Name = alias
			all = append(all, e)
			if strings.Contains(alias, "_item_no_check") {
				e.Name = strings.Replace(alias, "_item_no_check", "item_no_check", 1)
				all = append(all, e)
			}
		}
	}
	data, err := marshalLines(all)
	if err != nil {
		panic(err)
	}
	data = append(append([]byte("{\"macros\":"), data...), '}', '\n')
	if err := os.WriteFile("vocab/testdata/decomp_macros.json", data, 0644); err != nil {
		panic(err)
	}
	if err := generateConstants(root); err != nil {
		panic(err)
	}
	if err := generateReference(append(script, movement...)); err != nil {
		panic(err)
	}
	fmt.Printf("generated %d script and %d movement macros\n", len(script), len(movement))
}

func generateReference(entries []entry) error {
	if err := os.MkdirAll("docs", 0755); err != nil {
		return err
	}
	var b strings.Builder
	b.WriteString("# Command reference\n\nGenerated from the public decomp's `script.inc` and `movement.inc` macro declarations. Call these names in a `.poryz` script or movement block. Optional parameters show their default. Use `poryz constants PREFIX` to find flag, variable, item, species, move, and map names.\n\n| Command | Parameters | Area |\n| --- | --- | --- |\n")
	for _, e := range entries {
		area := "Script"
		if e.Movement {
			area = "Movement"
		}
		fmt.Fprintf(&b, "| `%s` | `%s` | %s |\n", e.Name, strings.Join(e.Params, ", "), area)
	}
	return os.WriteFile("docs/command-reference.md", []byte(b.String()), 0644)
}

func generateConstants(root string) error {
	type constant struct {
		Name  string `json:"name"`
		Value string `json:"value"`
		Kind  string `json:"kind"`
	}
	var all []constant
	for _, kind := range []string{"flags", "vars", "items", "species", "moves", "maps"} {
		file, err := os.Open(filepath.Join(root, "include/constants", kind+".h"))
		if err != nil {
			return err
		}
		s := bufio.NewScanner(file)
		for s.Scan() {
			fields := strings.Fields(strings.SplitN(s.Text(), "//", 2)[0])
			if len(fields) < 3 || fields[0] != "#define" || strings.Contains(fields[1], "(") {
				continue
			}
			all = append(all, constant{Name: fields[1], Value: strings.Join(fields[2:], " "), Kind: kind})
		}
		if err := s.Err(); err != nil {
			file.Close()
			return err
		}
		file.Close()
	}
	sort.Slice(all, func(i, j int) bool { return all[i].Name < all[j].Name })
	data, err := marshalLines(all)
	if err != nil {
		return err
	}
	return os.WriteFile("vocab/testdata/constants.json", data, 0644)
}

func marshalLines[T any](items []T) ([]byte, error) {
	var b bytes.Buffer
	b.WriteString("[\n")
	for i, item := range items {
		row, err := json.Marshal(item)
		if err != nil {
			return nil, err
		}
		b.Write(row)
		if i+1 < len(items) {
			b.WriteByte(',')
		}
		b.WriteByte('\n')
	}
	b.WriteString("]\n")
	return b.Bytes(), nil
}

func snake(name string) string {
	name = strings.Replace(name, "GoTo", "Goto", 1)
	r := []rune(name)
	var b strings.Builder
	for i, c := range r {
		if unicode.IsUpper(c) && i > 0 && (unicode.IsLower(r[i-1]) || unicode.IsDigit(r[i-1]) || (i+1 < len(r) && unicode.IsLower(r[i+1]) && unicode.IsUpper(r[i-1]))) {
			b.WriteByte('_')
		}
		b.WriteRune(unicode.ToLower(c))
	}
	return b.String()
}
