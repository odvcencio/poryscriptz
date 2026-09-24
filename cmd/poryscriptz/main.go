// Command poryscriptz compiles a .poryz event-script source file into
// HeartGold scr_seq macro-asm.
//
// Usage:
//
//	poryscriptz [--scrcmd <path>] [-o <out.s>] <file.poryz>
//
// --scrcmd defaults to "vocab/testdata/scrcmd.json" (relative to the working
// directory), which works when invoked from the repository root.
// -o writes output to a file; omit to write to stdout.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	poryscriptz "github.com/odvcencio/poryscriptz"
	"github.com/odvcencio/poryscriptz/target"
	"github.com/odvcencio/poryscriptz/vocab"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

// run is the testable entry point.  args is os.Args[1:].
// Returns 0 on success, 1 on error.
func run(args []string, stdout, stderr io.Writer) int {
	if len(args) > 0 {
		switch args[0] {
		case "help", "--help", "-h":
			fmt.Fprint(stdout, "poryz 0.3.0 — event scripts for HGSS\n\nCommands:\n  poryz decompile [-o out.poryz] file.s\n  poryz compile [-o out.s] file.poryz\n  poryz check file.poryz...\n  poryz fmt [-w] file.poryz...\n  poryz commands [PREFIX]\n  poryz constants [PREFIX]\n  poryz lsp\n\nRun poryz commands to browse the macro vocabulary.\n")
			return 0
		case "version", "--version":
			fmt.Fprintln(stdout, "poryz 0.3.0")
			return 0
		case "compile":
			return runCompile(args[1:], stdout, stderr)
		case "decompile":
			return runDecompile(args[1:], stdout, stderr)
		case "check":
			return runCheck(args[1:], stdout, stderr)
		case "fmt":
			return runFmt(args[1:], stdout, stderr)
		case "constants":
			return runConstants(args[1:], stdout, stderr)
		case "commands":
			return runCommands(args[1:], stdout, stderr)
		case "lsp":
			return runLSP(os.Stdin, stdout, stderr)
		}
	}
	return runCompile(args, stdout, stderr)
}

func runCompile(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("poryscriptz", flag.ContinueOnError)
	fs.SetOutput(stderr)
	scrcmd := fs.String("scrcmd", "", "path to scrcmd.json vocabulary file (default: built-in)")
	out := fs.String("o", "", "output file (default: stdout)")

	if err := fs.Parse(args); err != nil {
		return 1
	}
	if fs.NArg() != 1 {
		fmt.Fprintf(stderr, "usage: poryz compile [--scrcmd <path>] [-o <out.s>] <file.poryz>\n")
		return 1
	}
	inputFile := fs.Arg(0)

	tgt, err := target.HGSS(*scrcmd)
	if err != nil {
		fmt.Fprintf(stderr, "poryscriptz: load vocabulary: %v\n", err)
		return 1
	}

	src, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Fprintf(stderr, "poryscriptz: read %s: %v\n", inputFile, err)
		return 1
	}

	asm, diags, err := poryscriptz.Compile(src, tgt)
	if err != nil {
		fmt.Fprintf(stderr, "%s:%v; check braces, commas, and call syntax\n", inputFile, err)
		return 1
	}
	if len(diags) > 0 {
		for _, d := range diags {
			fmt.Fprintf(stderr, "%s:%d:%d: %s\n", inputFile, d.Line, d.Col, d.Msg)
		}
		return 1
	}

	if *out == "" {
		fmt.Fprint(stdout, asm)
		return 0
	}
	if err := os.WriteFile(*out, []byte(asm), 0644); err != nil {
		fmt.Fprintf(stderr, "poryscriptz: write %s: %v\n", *out, err)
		return 1
	}
	return 0
}

func runDecompile(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("decompile", flag.ContinueOnError)
	fs.SetOutput(stderr)
	scrcmd := fs.String("scrcmd", "", "path to scrcmd.json vocabulary file (default: built-in)")
	out := fs.String("o", "", "output .poryz file (default: stdout)")
	if err := fs.Parse(args); err != nil {
		return 1
	}
	if fs.NArg() != 1 {
		fmt.Fprintln(stderr, "usage: poryz decompile [--scrcmd path] [-o out.poryz] <scr_seq.s>")
		return 1
	}
	input := fs.Arg(0)
	tgt, err := target.HGSS(*scrcmd)
	if err != nil {
		fmt.Fprintf(stderr, "poryz: load vocabulary: %v\n", err)
		return 1
	}
	src, err := os.ReadFile(input)
	if err != nil {
		fmt.Fprintf(stderr, "%s: %v\n", input, err)
		return 1
	}
	code, err := poryscriptz.Decompile(src, tgt)
	if err != nil {
		fmt.Fprintf(stderr, "%s:%v\n", input, err)
		return 1
	}
	if *out == "" {
		fmt.Fprint(stdout, code)
		return 0
	}
	if err := os.WriteFile(*out, []byte(code), 0644); err != nil {
		fmt.Fprintf(stderr, "%s: %v\n", *out, err)
		return 1
	}
	return 0
}

func runCheck(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("check", flag.ContinueOnError)
	fs.SetOutput(stderr)
	scrcmd := fs.String("scrcmd", "", "path to scrcmd.json vocabulary file (default: built-in)")
	if err := fs.Parse(args); err != nil {
		return 1
	}
	if fs.NArg() == 0 {
		fmt.Fprintln(stderr, "usage: poryz check [--scrcmd path] <file.poryz>...")
		return 1
	}
	tgt, err := target.HGSS(*scrcmd)
	if err != nil {
		fmt.Fprintf(stderr, "poryz: load vocabulary: %v\n", err)
		return 1
	}
	failed := false
	for _, path := range fs.Args() {
		src, err := os.ReadFile(path)
		if err != nil {
			fmt.Fprintf(stderr, "%s: %v\n", path, err)
			failed = true
			continue
		}
		_, diags, err := poryscriptz.Compile(src, tgt)
		if err != nil {
			fmt.Fprintf(stderr, "%s:%v\n", path, err)
			failed = true
			continue
		}
		for _, d := range diags {
			fmt.Fprintf(stderr, "%s:%d:%d: %s\n", path, d.Line, d.Col, d.Msg)
			failed = true
		}
	}
	if failed {
		return 1
	}
	fmt.Fprintf(stdout, "checked %d file(s)\n", fs.NArg())
	return 0
}

func runFmt(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("fmt", flag.ContinueOnError)
	fs.SetOutput(stderr)
	write := fs.Bool("w", false, "write result to each file")
	if err := fs.Parse(args); err != nil {
		return 1
	}
	if fs.NArg() == 0 {
		fmt.Fprintln(stderr, "usage: poryz fmt [-w] <file.poryz>...")
		return 1
	}
	failed := false
	for _, file := range fs.Args() {
		src, err := os.ReadFile(file)
		if err != nil {
			fmt.Fprintf(stderr, "%s: %v\n", file, err)
			failed = true
			continue
		}
		formatted, err := poryscriptz.Format(src)
		if err != nil {
			fmt.Fprintf(stderr, "%s:%v; check syntax before formatting\n", file, err)
			failed = true
			continue
		}
		formatted = append(formatted, '\n')
		if *write {
			err = os.WriteFile(file, formatted, 0644)
		} else {
			_, err = stdout.Write(formatted)
		}
		if err != nil {
			fmt.Fprintf(stderr, "%s: %v\n", file, err)
			failed = true
		}
	}
	if failed {
		return 1
	}
	return 0
}

func runConstants(args []string, stdout, stderr io.Writer) int {
	if len(args) > 1 {
		fmt.Fprintln(stderr, "usage: poryz constants [PREFIX]")
		return 1
	}
	prefix := ""
	if len(args) == 1 {
		prefix = args[0]
	}
	all, err := vocab.Constants(prefix)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	for _, c := range all {
		fmt.Fprintf(stdout, "%s = %s\n", c.Name, c.Value)
	}
	if len(all) == 0 {
		fmt.Fprintf(stderr, "no constants match %q; try a shorter prefix\n", prefix)
		return 1
	}
	return 0
}

func runCommands(args []string, stdout, stderr io.Writer) int {
	if len(args) > 1 {
		fmt.Fprintln(stderr, "usage: poryz commands [PREFIX]")
		return 1
	}
	prefix := ""
	if len(args) == 1 {
		prefix = strings.ToLower(args[0])
	}
	tgt, err := target.HGSS("")
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	entries := tgt.Macros().Entries()
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name < entries[j].Name })
	count := 0
	for _, e := range entries {
		if prefix == "" && (e.Name == "" || e.Name[0] < 'A' || e.Name[0] > 'Z') {
			continue
		}
		if !strings.HasPrefix(strings.ToLower(e.Name), prefix) {
			continue
		}
		params := e.Params
		if len(params) == 0 && len(e.Args) > 0 {
			params = make([]string, len(e.Args))
			for i := range params {
				params[i] = "arg"
			}
		}
		fmt.Fprintf(stdout, "%s(%s)\n", e.Name, strings.Join(params, ", "))
		count++
	}
	if count == 0 {
		fmt.Fprintf(stderr, "no commands match %q; try a shorter prefix\n", prefix)
		return 1
	}
	return 0
}
