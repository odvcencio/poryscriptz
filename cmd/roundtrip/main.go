// roundtrip assembles original and decompiled scripts and compares binary bytes.
// It needs a local MWAS installation and the decomp's generated headers.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	poryz "github.com/odvcencio/poryscriptz"
	"github.com/odvcencio/poryscriptz/target"
)

type result struct{ file, class, detail string }

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr, check))
}

type checker func(file, root, mwas, headers string, tgt target.GameTarget) result

func run(args []string, stdout, stderr io.Writer, checkFile checker) int {
	flags := flag.NewFlagSet("roundtrip", flag.ContinueOnError)
	flags.SetOutput(stderr)
	decomp := flags.String("decomp", "", "public decomp root")
	hack := flags.String("hack", "", "hack scr_seq directory")
	hackRoot := flags.String("hack-root", "", "decomp root with hack macros and headers")
	mwas := flags.String("mwas", "", "mwasmarm.exe path")
	headers := flags.String("headers", "", "extra generated header files root")
	vocabulary := flags.String("vocab", "vocab/testdata/scrcmd.json", "scrcmd.json path")
	workers := flags.Int("workers", 6, "parallel assembler processes")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if *decomp == "" || *mwas == "" {
		fmt.Fprintln(stderr, "usage: roundtrip -decomp ROOT -mwas EXE [-headers FILES_ROOT] [-hack SCRIPTS -hack-root ROOT]")
		return 2
	}
	if *workers < 1 || flags.NArg() != 0 {
		fmt.Fprintln(stderr, "roundtrip: workers must be >= 1 and positional arguments are not supported")
		return 2
	}
	tgt, err := target.HGSS(*vocabulary)
	if err != nil {
		fmt.Fprintf(stderr, "roundtrip: load vocabulary: %v\n", err)
		return 2
	}
	failed := false
	for _, set := range []struct{ name, folder, root string }{{"pret", filepath.Join(*decomp, "files/fielddata/script/scr_seq"), *decomp}, {"hack", *hack, *hackRoot}} {
		if set.folder == "" {
			continue
		}
		if set.root == "" {
			set.root = *decomp
		}
		files, err := filepath.Glob(filepath.Join(set.folder, "*.s"))
		if err != nil {
			fmt.Fprintf(stderr, "%s: list corpus: %v\n", set.name, err)
			failed = true
			continue
		}
		if len(files) == 0 {
			fmt.Fprintf(stderr, "%s: empty corpus: expected .s files in %s\n", set.name, set.folder)
			failed = true
			continue
		}
		jobs := make(chan string)
		results := make(chan result, len(files))
		var wg sync.WaitGroup
		for i := 0; i < *workers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for file := range jobs {
					results <- checkFile(file, set.root, *mwas, *headers, tgt)
				}
			}()
		}
		for _, file := range files {
			jobs <- file
		}
		close(jobs)
		wg.Wait()
		close(results)
		classes := map[string][]result{}
		pass := 0
		for r := range results {
			if r.class == "" {
				pass++
			} else {
				classes[r.class] = append(classes[r.class], r)
			}
		}
		if pass != len(files) {
			failed = true
		}
		fmt.Fprintf(stdout, "%s: %d/%d byte-identical\n", set.name, pass, len(files))
		keys := make([]string, 0, len(classes))
		for k := range classes {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			sort.Slice(classes[k], func(i, j int) bool { return classes[k][i].file < classes[k][j].file })
			fmt.Fprintf(stdout, "  %s: %d\n", k, len(classes[k]))
			for _, r := range classes[k] {
				fmt.Fprintf(stdout, "    %s: %s\n", filepath.Base(r.file), r.detail)
			}
		}
	}
	if failed {
		return 1
	}
	return 0
}

func check(file, root, mwas, headers string, tgt target.GameTarget) result {
	r := result{file: file}
	src, err := os.ReadFile(file)
	if err != nil {
		r.class = "read"
		r.detail = err.Error()
		return r
	}
	z, err := poryz.Decompile(src, tgt)
	if err != nil {
		r.class = "decompile"
		r.detail = err.Error()
		return r
	}
	asm, diags, err := poryz.Compile([]byte(z), tgt)
	if err != nil {
		r.class = "compile"
		r.detail = err.Error()
		return r
	}
	if len(diags) > 0 {
		r.class = "compile"
		r.detail = diags[0].Msg
		return r
	}
	tmp, err := os.MkdirTemp("", "poryz-roundtrip-")
	if err != nil {
		r.class = "temp"
		r.detail = err.Error()
		return r
	}
	defer os.RemoveAll(tmp)
	gen := filepath.Join(tmp, "generated.s")
	if err = os.WriteFile(gen, []byte(asm), 0644); err != nil {
		r.class = "write"
		r.detail = err.Error()
		return r
	}
	oldBytes, err := assemble(file, tmp, "original", root, mwas, headers)
	if err != nil {
		r.class = "original assembly"
		r.detail = err.Error()
		return r
	}
	newBytes, err := assemble(gen, tmp, "generated", root, mwas, headers)
	if err != nil {
		r.class = "generated assembly"
		r.detail = err.Error()
		return r
	}
	if !bytes.Equal(oldBytes, newBytes) {
		r.class = "byte mismatch"
		r.detail = fmt.Sprintf("original %d bytes, generated %d bytes", len(oldBytes), len(newBytes))
	}
	return r
}

func assemble(source, tmp, name, root, mwas, headers string) ([]byte, error) {
	object := filepath.Join(tmp, name+".o")
	binary := filepath.Join(tmp, name+".bin")
	args := []string{mwas, "-DHEARTGOLD", "-DGAME_REMASTER=0", "-DENGLISH", "-DPM_KEEP_ASSERTS", "-DSDK_ARM9", "-DSDK_CODE_ARM", "-DSDK_FINALROM", "-DPM_ASM", "-DSDK_ASM", "-proc", "arm5te", "-g", "-gccinc"}
	for _, path := range []string{".", "./include", "./asm/include", "./files", headers, "./lib/asm/include", "./lib/NitroDWC/asm/include", "./lib/MSL_C/asm/include", "./lib/NitroSDK/asm/include", "./lib/syscall/asm/include", "./asm", "./files/msgdata"} {
		if path != "" {
			args = append(args, "-i", path)
		}
	}
	args = append(args, "-I./lib/include", "-o", object, source)
	cmd := exec.Command("wine", args...)
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "WINEARCH=win32")
	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("mwas: %s: %w", strings.TrimSpace(string(out)), err)
	}
	cmd = exec.Command("arm-none-eabi-objcopy", "-O", "binary", "--file-alignment", "4", object, binary)
	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("objcopy: %s: %w", strings.TrimSpace(string(out)), err)
	}
	return os.ReadFile(binary)
}
