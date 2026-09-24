// roundtrip assembles original and decompiled scripts and compares binary bytes.
// It needs a local MWAS installation and the decomp's generated headers.
package main

import (
	"bytes"
	"flag"
	"fmt"
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
	decomp := flag.String("decomp", "", "public decomp root")
	hack := flag.String("hack", "", "hack scr_seq directory")
	hackRoot := flag.String("hack-root", "", "decomp root with hack macros and headers")
	mwas := flag.String("mwas", "", "mwasmarm.exe path")
	headers := flag.String("headers", "", "extra generated header files root")
	vocabulary := flag.String("vocab", "vocab/testdata/scrcmd.json", "scrcmd.json path")
	workers := flag.Int("workers", 6, "parallel assembler processes")
	flag.Parse()
	if *decomp == "" || *mwas == "" {
		fmt.Fprintln(os.Stderr, "usage: roundtrip -decomp ROOT -mwas EXE [-headers FILES_ROOT] [-hack SCRIPTS -hack-root ROOT]")
		os.Exit(2)
	}
	tgt, err := target.HGSS(*vocabulary)
	if err != nil {
		panic(err)
	}
	for _, set := range []struct{ name, folder, root string }{{"pret", filepath.Join(*decomp, "files/fielddata/script/scr_seq"), *decomp}, {"hack", *hack, *hackRoot}} {
		if set.folder == "" {
			continue
		}
		if set.root == "" {
			set.root = *decomp
		}
		files, err := filepath.Glob(filepath.Join(set.folder, "*.s"))
		if err != nil {
			panic(err)
		}
		jobs := make(chan string)
		results := make(chan result, len(files))
		var wg sync.WaitGroup
		for i := 0; i < *workers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for file := range jobs {
					results <- check(file, set.root, *mwas, *headers, tgt)
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
		fmt.Printf("%s: %d/%d byte-identical\n", set.name, pass, len(files))
		keys := make([]string, 0, len(classes))
		for k := range classes {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			sort.Slice(classes[k], func(i, j int) bool { return classes[k][i].file < classes[k][j].file })
			fmt.Printf("  %s: %d\n", k, len(classes[k]))
			for _, r := range classes[k] {
				fmt.Printf("    %s: %s\n", filepath.Base(r.file), r.detail)
			}
		}
	}
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
