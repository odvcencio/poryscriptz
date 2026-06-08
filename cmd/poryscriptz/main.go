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

	poryscriptz "m31labs.dev/poryscriptz"
	"m31labs.dev/poryscriptz/target"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

// run is the testable entry point.  args is os.Args[1:].
// Returns 0 on success, 1 on error.
func run(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("poryscriptz", flag.ContinueOnError)
	fs.SetOutput(stderr)
	scrcmd := fs.String("scrcmd", "vocab/testdata/scrcmd.json", "path to scrcmd.json vocabulary file")
	out := fs.String("o", "", "output file (default: stdout)")

	if err := fs.Parse(args); err != nil {
		return 1
	}
	if fs.NArg() != 1 {
		fmt.Fprintf(stderr, "usage: poryscriptz [--scrcmd <path>] [-o <out.s>] <file.poryz>\n")
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
		fmt.Fprintf(stderr, "poryscriptz: %v\n", err)
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
