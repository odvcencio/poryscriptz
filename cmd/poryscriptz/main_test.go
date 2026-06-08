package main

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// repoRoot returns the repository root (two levels above this file's directory).
func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	// file = .../cmd/poryscriptz/main_test.go → go up two dirs
	return filepath.Join(filepath.Dir(file), "..", "..")
}

func TestCLIRoute30(t *testing.T) {
	root := repoRoot(t)
	golden, err := os.ReadFile(filepath.Join(root, "testdata", "route30.golden"))
	if err != nil {
		t.Fatalf("read golden: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := run([]string{
		"--scrcmd", filepath.Join(root, "vocab", "testdata", "scrcmd.json"),
		filepath.Join(root, "testdata", "route30.poryz"),
	}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("run returned %d; stderr: %s", code, stderr.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("unexpected stderr: %s", stderr.String())
	}
	if stdout.String() != string(golden) {
		t.Fatalf("stdout mismatch:\n--- got ---\n%s\n--- want ---\n%s", stdout.String(), string(golden))
	}
}

func TestCLIUnknownCommand(t *testing.T) {
	root := repoRoot(t)

	// Write a temp file with an unknown command.
	tmp, err := os.CreateTemp(t.TempDir(), "bad_*.poryz")
	if err != nil {
		t.Fatalf("create temp: %v", err)
	}
	_, _ = tmp.WriteString("package m\nscript S {\n\tnot_a_real_command()\n}\n")
	_ = tmp.Close()

	var stdout, stderr bytes.Buffer
	code := run([]string{
		"--scrcmd", filepath.Join(root, "vocab", "testdata", "scrcmd.json"),
		tmp.Name(),
	}, &stdout, &stderr)

	if code == 0 {
		t.Fatal("expected non-zero exit for unknown command, got 0")
	}
	errOut := stderr.String()
	// Diagnostic must be in the form <file>:<line>:<col>: <msg>
	if !strings.Contains(errOut, tmp.Name()) {
		t.Errorf("stderr should contain filename %q, got: %s", tmp.Name(), errOut)
	}
	if !strings.Contains(errOut, "not in") {
		t.Errorf("stderr should contain 'not in', got: %s", errOut)
	}
	// Must contain line:col pattern (non-zero line)
	if !strings.Contains(errOut, ":3:") {
		t.Errorf("stderr should contain ':3:' (line 3), got: %s", errOut)
	}
}
