package poryscriptz

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/odvcencio/poryscriptz/target"
)

var asmLabel = regexp.MustCompile(`^([A-Za-z_][A-Za-z_0-9]*):(?:\s*;.*)?$`)
var asmDefine = regexp.MustCompile(`^#define\s+([A-Za-z_][A-Za-z_0-9]*)\s+(.+)$`)

type asmLine struct {
	text string
	line int
}
type asmBlock struct {
	name     string
	lines    []asmLine
	movement bool
	source   string
}

// Decompile converts a scr_seq macro assembly file to editable poryz source.
// It rejects lines whose bytes cannot be reproduced by the compiler.
func Decompile(src []byte, tgt target.GameTarget) (string, error) {
	var includes, entries []string
	var blocks []asmBlock
	var head []asmLine
	var defines = map[string]string{}
	var out strings.Builder
	var initHeader, sawScrDefEnd, sawBody bool
	var tableEnds int
	var pretDialect bool
	var pendingAlign bool
	current := -1
	scanner := bufio.NewScanner(strings.NewReader(string(src)))
	scanner.Buffer(make([]byte, 4096), 2*1024*1024)
	for lineNo := 1; scanner.Scan(); lineNo++ {
		raw := strings.TrimSpace(scanner.Text())
		if raw == "" || strings.HasPrefix(raw, ";") || strings.HasPrefix(raw, "//") {
			continue
		}
		if m := asmDefine.FindStringSubmatch(raw); m != nil {
			defines[m[1]] = strings.TrimSpace(m[2])
			continue
		}
		if strings.HasPrefix(raw, "#include ") {
			value := strings.TrimSpace(strings.TrimPrefix(raw, "#include "))
			if value == `"constants/scrcmd.h"` {
				continue
			}
			header, err := strconv.Unquote(value)
			if err != nil || !strings.HasSuffix(header, ".h") {
				return "", fmt.Errorf("%d:1: unsupported include %q; use a quoted .h path", lineNo, value)
			}
			includes = append(includes, header)
			continue
		}
		if raw == `.include "asm/macros/script.inc"` || raw == ".rodata" || raw == ".option alignment off" {
			continue
		}
		if strings.HasPrefix(raw, `.include `) {
			path, err := strconv.Unquote(strings.TrimSpace(strings.TrimPrefix(raw, `.include `)))
			if err != nil || !strings.HasSuffix(path, ".s") {
				return "", fmt.Errorf("%d:1: unsupported source include; use a relative .s path", lineNo)
			}
			if pendingAlign {
				blocks = append(blocks, asmBlock{name: "@align"})
				pendingAlign = false
			}
			blocks = append(blocks, asmBlock{name: "@source", source: path})
			current = -1
			continue
		}
		if strings.HasPrefix(raw, ".balign") {
			if raw != ".balign 4, 0" {
				return "", fmt.Errorf("%d:1: unsupported alignment %q; use .balign 4, 0", lineNo, raw)
			}
			pendingAlign = true
			continue
		}
		if m := asmLabel.FindStringSubmatch(raw); m != nil {
			if pendingAlign {
				blocks = append(blocks, asmBlock{name: "@align"})
				pendingAlign = false
			}
			blocks = append(blocks, asmBlock{name: m[1]})
			current = len(blocks) - 1
			sawBody = true
			continue
		}
		cmd, args, ok := parseAsmCall(raw)
		if !ok {
			return "", fmt.Errorf("%d:1: unsupported assembly %q; use a script.inc macro or .balign 4, 0", lineNo, raw)
		}
		if cmd == "ScrDef" || cmd == "scrdef" {
			if cmd == "ScrDef" {
				pretDialect = true
			}
			if sawBody || len(args) != 1 {
				return "", fmt.Errorf("%d:1: ScrDef must have one label before script bodies", lineNo)
			}
			entries = append(entries, args[0])
			continue
		}
		if cmd == "ScrDefEnd" || cmd == "scrdef_end" {
			sawScrDefEnd = true
			tableEnds++
			if cmd == "ScrDefEnd" {
				pretDialect = true
			}
			continue
		}
		if strings.HasPrefix(cmd, "InitScript") {
			pretDialect = true
		}
		if macro, ok := tgt.Macros().ByName(cmd); ok {
			if len(args) < macro.MinArgs || len(args) > len(macro.Args) {
				return "", fmt.Errorf("%d:2: macro %s expects %d to %d args, got %d; check the command reference", lineNo, cmd, macro.MinArgs, len(macro.Args), len(args))
			}
		} else if command, ok := tgt.Vocabulary().ByName(cmd); ok {
			if len(args) != len(command.Args) {
				return "", fmt.Errorf("%d:2: command %s expects %d args, got %d; check the command reference", lineNo, cmd, len(command.Args), len(args))
			}
		} else {
			return "", fmt.Errorf("%d:2: unknown macro %q; regenerate the vocabulary from this decomp's asm/macros/script.inc", lineNo, cmd)
		}
		for i, arg := range args {
			for name, value := range defines {
				arg = regexp.MustCompile(`\b`+regexp.QuoteMeta(name)+`\b`).ReplaceAllString(arg, value)
			}
			args[i] = arg
		}
		if alias := map[string]string{"goto": "AsmGoto", "return": "AsmReturn", "switch": "AsmSwitch", "case": "AsmCase"}[cmd]; alias != "" {
			cmd = alias
		}
		call := asmLine{text: cmd + "(" + strings.Join(args, ", ") + ")", line: lineNo}
		if current < 0 {
			initHeader = true
			head = append(head, call)
		} else {
			blocks[current].lines = append(blocks[current].lines, call)
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	if pendingAlign {
		blocks = append(blocks, asmBlock{name: "@align"})
	}
	if !sawScrDefEnd {
		initHeader = true
	}
	if initHeader && sawScrDefEnd {
		return "", fmt.Errorf("1:1: mixed script table and init header; split into separate .s files")
	}
	out.WriteString("package m\n")
	if pretDialect {
		out.WriteString("dialect pret\n")
	} else {
		out.WriteString("dialect legacy\n")
	}
	for _, h := range includes {
		fmt.Fprintf(&out, "include %q\n", h)
	}
	if initHeader {
		out.WriteString("\nheader init {\n")
		for _, l := range head {
			fmt.Fprintf(&out, "    %s\n", l.text)
		}
		out.WriteString("}\n")
	} else if len(entries) > 0 {
		out.WriteString("\nentries(")
		out.WriteString(strings.Join(entries, ", "))
		out.WriteString(")\n")
	}
	for i := 1; i < tableEnds; i++ {
		out.WriteString("table_end()\n")
	}
	entrySet := map[string]bool{}
	for _, e := range entries {
		entrySet[e] = true
	}
	for _, b := range blocks {
		out.WriteByte('\n')
		if b.name == "@align" {
			out.WriteString("align4()\n")
			continue
		}
		if b.name == "@source" {
			fmt.Fprintf(&out, "source %q\n", b.source)
			continue
		}
		kind := "label"
		if entrySet[b.name] {
			kind = "script"
		}
		fmt.Fprintf(&out, "%s %s {\n", kind, b.name)
		for _, l := range b.lines {
			fmt.Fprintf(&out, "    %s\n", l.text)
		}
		out.WriteString("}\n")
	}
	code := out.String()
	if _, diags, err := Compile([]byte(code), tgt); err != nil {
		return "", fmt.Errorf("1:1: generated source has syntax error: %v; check assembly operands", err)
	} else if len(diags) > 0 {
		return "", fmt.Errorf("%d:%d: generated source: %s", diags[0].Line, diags[0].Col, diags[0].Msg)
	}
	return code, nil
}

func parseAsmCall(line string) (string, []string, bool) {
	line = strings.SplitN(line, ";", 2)[0]
	line = strings.TrimSpace(line)
	if line == "" {
		return "", nil, false
	}
	name, rest, _ := strings.Cut(line, " ")
	if !regexp.MustCompile(`^[A-Za-z_][A-Za-z_0-9]*$`).MatchString(name) {
		return "", nil, false
	}
	rest = strings.TrimSpace(rest)
	if rest == "" {
		return name, nil, true
	}
	var args []string
	depth, start := 0, 0
	for i, c := range rest {
		switch c {
		case '(':
			depth++
		case ')':
			depth--
		case ',':
			if depth == 0 {
				args = append(args, strings.TrimSpace(rest[start:i]))
				start = i + 1
			}
		}
	}
	args = append(args, strings.TrimSpace(rest[start:]))
	return name, args, true
}
