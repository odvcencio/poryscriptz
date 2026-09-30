package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/odvcencio/poryscriptz/target"
)

func corpus(t *testing.T) (string, []string) {
	t.Helper()
	root := t.TempDir()
	folder := filepath.Join(root, "files", "fielddata", "script", "scr_seq")
	if err := os.MkdirAll(folder, 0755); err != nil {
		t.Fatal(err)
	}
	args := []string{"-decomp", root, "-mwas", "unused.exe", "-vocab", "../../vocab/testdata/scrcmd.json"}
	return folder, args
}

func TestRoundtripEmptyRequiredOrRequestedCorpusFails(t *testing.T) {
	folder, args := corpus(t)
	var out, errOut bytes.Buffer
	if code := run(args, &out, &errOut, check); code != 1 || !strings.Contains(errOut.String(), "pret: empty corpus") {
		t.Fatalf("empty public corpus: code=%d stdout=%s stderr=%s", code, &out, &errOut)
	}
	if err := os.WriteFile(filepath.Join(folder, "one.s"), nil, 0644); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	errOut.Reset()
	pass := func(file, _, _, _ string, _ target.GameTarget) result { return result{file: file} }
	args = append(args, "-hack", t.TempDir())
	if code := run(args, &out, &errOut, pass); code != 1 || !strings.Contains(errOut.String(), "hack: empty corpus") {
		t.Fatalf("empty requested edited corpus: code=%d stderr=%s", code, &errOut)
	}
}

func TestRoundtripEveryFailureClassMakesGateFail(t *testing.T) {
	folder, args := corpus(t)
	for _, name := range []string{"z.s", "a.s"} {
		if err := os.WriteFile(filepath.Join(folder, name), nil, 0644); err != nil {
			t.Fatal(err)
		}
	}
	for _, class := range []string{"read", "decompile", "compile", "temp", "write", "original assembly", "generated assembly", "byte mismatch"} {
		t.Run(class, func(t *testing.T) {
			var out, errOut bytes.Buffer
			fail := func(file, _, _, _ string, _ target.GameTarget) result {
				return result{file: file, class: class, detail: "fixture failure"}
			}
			if code := run(args, &out, &errOut, fail); code != 1 || !strings.Contains(out.String(), "0/2 byte-identical") || !strings.Contains(out.String(), class+": 2") {
				t.Fatalf("code=%d stdout=%s stderr=%s", code, &out, &errOut)
			}
			if strings.Index(out.String(), "a.s:") > strings.Index(out.String(), "z.s:") {
				t.Fatal("failure report is not deterministically sorted")
			}
		})
	}
}

func TestRoundtripPassAndInvalidWorkers(t *testing.T) {
	folder, args := corpus(t)
	if err := os.WriteFile(filepath.Join(folder, "one.s"), nil, 0644); err != nil {
		t.Fatal(err)
	}
	pass := func(file, _, _, _ string, _ target.GameTarget) result { return result{file: file} }
	var out, errOut bytes.Buffer
	if code := run(args, &out, &errOut, pass); code != 0 || !strings.Contains(out.String(), "1/1 byte-identical") {
		t.Fatalf("valid corpus: code=%d output=%s", code, &out)
	}
	for _, workers := range []string{"0", "-1"} {
		if code := run(append(args, "-workers", workers), &out, &errOut, pass); code != 2 {
			t.Fatalf("workers=%s must fail before starting jobs: code=%d", workers, code)
		}
	}
}
