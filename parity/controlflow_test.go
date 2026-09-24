//go:build decompparity

// Package parity — control-flow assembles-clean gate (PCF-T3).
//
// For each v0.2 control-flow construct (ifelse, elseif, while, switch) and the
// v0.1 flag-gate / var-gate constructs, this test:
//
//  1. Compiles the corresponding .poryz fixture via the poryscriptZ compiler.
//  2. Writes the emitted .s to a scratch temp dir.
//  3. Runs the Metrowerks MWAS assembler (under Wine) with the same flags the
//     decomp Makefile uses — the same command that TestCompareScrSeq0027 uses.
//  4. Asserts exit 0 / no stderr error, proving the emitted bytecode + labels +
//     macros are valid for the HGSS assembler.
//
// The test does NOT check SHA1 byte-identity for these synthetic scripts (that
// is done only for scr_seq_0027 in parity_test.go).  Assembles-clean is the
// primary deliverable for PCF-T3.
//
// # Environment requirements (same as TestCompareScrSeq0027)
//
//   - Wine prefix: $POKEHG_WINEPREFIX or $WINEPREFIX, defaulting to ~/.wine-pokehg-build
//   - Metrowerks MWAS: $POKEHG_ROOT/tools/mwccarm/2.0/sp2p2/mwasmarm.exe
//   - $POKEHG_ROOT: path to the pokeheartgold decomp checkout
//
// # Running manually
//
//	go test -tags decompparity ./parity/ -v -run TestControlFlowAssemblesClean
package parity

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	poryscriptz "github.com/odvcencio/poryscriptz"
	"github.com/odvcencio/poryscriptz/target"
)

// cfCase describes one control-flow fixture to compile and assemble.
type cfCase struct {
	name string // sub-test name, e.g. "ifelse"
	file string // filename under testdata/parity/controlflow/
}

// controlFlowCases lists all v0.2 control-flow constructs plus v0.1 gates.
var controlFlowCases = []cfCase{
	{name: "ifelse", file: "ifelse.poryz"},
	{name: "elseif", file: "elseif.poryz"},
	{name: "while", file: "while.poryz"},
	{name: "switch", file: "switch.poryz"},
	{name: "flaggate", file: "flaggate.poryz"},
	{name: "vargate", file: "vargate.poryz"},
}

// macrosCases lists v0.3 macro-vocabulary fixtures (MacrosAssembleClean gate).
var macrosCases = []cfCase{
	{name: "macros", file: "macros.poryz"},
	{name: "movement", file: "movement.poryz"},
	{name: "patterns", file: "patterns.poryz"},
}

// mwasBaseArgs are the MWAS flags shared by all assembler invocations.
// They mirror the decomp Makefile exactly (same as TestCompareScrSeq0027).
var mwasBaseArgs = []string{
	"-DHEARTGOLD", "-DGAME_REMASTER=0", "-DENGLISH",
	"-DPM_KEEP_ASSERTS", "-DSDK_ARM9", "-DSDK_CODE_ARM", "-DSDK_FINALROM",
	"-DPM_ASM", "-DSDK_ASM",
	"-proc", "arm5te",
	"-g", "-gccinc",
	"-i", ".", "-i", "./include",
	"-i", "./asm/include",
	"-i", "./files",
	"-i", "./lib/asm/include",
	"-i", "./lib/NitroDWC/asm/include",
	"-i", "./lib/MSL_C/asm/include",
	"-i", "./lib/NitroSDK/asm/include",
	"-i", "./lib/syscall/asm/include",
	"-i", "./asm",
	"-i", "./files/msgdata",
	"-I./lib/include",
}

// TestControlFlowAssemblesClean compiles each v0.2 control-flow .poryz
// fixture and asserts the emitted .s assembles cleanly under MWAS.
func TestControlFlowAssemblesClean(t *testing.T) {
	// --- prerequisite checks (same as TestCompareScrSeq0027) ---
	hgRoot := pokehgRoot(t)
	mwas := mwasPath(hgRoot)
	wp := winePrefix()

	for _, check := range []struct {
		label string
		path  string
	}{
		{"decomp root", hgRoot},
		{"mwasmarm.exe", mwas},
		{"wine prefix", wp},
	} {
		if _, err := os.Stat(check.path); err != nil {
			t.Skipf("decompparity: prerequisite not found (%s=%s): %v", check.label, check.path, err)
		}
	}
	if _, err := exec.LookPath("wine"); err != nil {
		t.Skip("decompparity: wine not in PATH")
	}

	// --- locate testdata/parity/controlflow/ ---
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("decompparity: runtime.Caller failed")
	}
	cfDir := filepath.Join(filepath.Dir(thisFile), "..", "testdata", "parity", "controlflow")

	// --- load HGSS target once ---
	gt, err := target.HGSS(scrcmdJSON(t))
	if err != nil {
		t.Fatalf("load HGSS target: %v", err)
	}

	// --- run one sub-test per fixture ---
	for _, tc := range controlFlowCases {
		tc := tc // capture
		t.Run(tc.name, func(t *testing.T) {
			runAssemblesClean(t, tc, cfDir, gt, hgRoot, mwas, wp)
		})
	}
}

// TestMacrosAssembleClean compiles each v0.3 macro-vocabulary .poryz fixture
// and asserts the emitted .s assembles cleanly under MWAS (exit 0).
// This is the v0.3 macro-vocabulary assembles-clean gate (analogous to
// TestControlFlowAssemblesClean for v0.2 control-flow constructs).
func TestMacrosAssembleClean(t *testing.T) {
	// --- prerequisite checks (same as TestCompareScrSeq0027) ---
	hgRoot := pokehgRoot(t)
	mwas := mwasPath(hgRoot)
	wp := winePrefix()

	for _, check := range []struct {
		label string
		path  string
	}{
		{"decomp root", hgRoot},
		{"mwasmarm.exe", mwas},
		{"wine prefix", wp},
	} {
		if _, err := os.Stat(check.path); err != nil {
			t.Skipf("decompparity: prerequisite not found (%s=%s): %v", check.label, check.path, err)
		}
	}
	if _, err := exec.LookPath("wine"); err != nil {
		t.Skip("decompparity: wine not in PATH")
	}

	// --- locate testdata/parity/controlflow/ ---
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("decompparity: runtime.Caller failed")
	}
	cfDir := filepath.Join(filepath.Dir(thisFile), "..", "testdata", "parity", "controlflow")

	// --- load HGSS target once ---
	gt, err := target.HGSS(scrcmdJSON(t))
	if err != nil {
		t.Fatalf("load HGSS target: %v", err)
	}

	// --- run one sub-test per macro fixture ---
	for _, tc := range macrosCases {
		tc := tc // capture
		t.Run(tc.name, func(t *testing.T) {
			runAssemblesClean(t, tc, cfDir, gt, hgRoot, mwas, wp)
		})
	}
}

// runAssemblesClean is the per-case worker: compile .poryz → write .s → assemble.
func runAssemblesClean(
	t *testing.T,
	tc cfCase,
	cfDir string,
	gt target.GameTarget,
	hgRoot, mwas, wp string,
) {
	t.Helper()

	// --- compile .poryz → .s ---
	poryzPath := filepath.Join(cfDir, tc.file)
	src, err := os.ReadFile(poryzPath)
	if err != nil {
		t.Fatalf("read %s: %v", tc.file, err)
	}

	asm, diags, compErr := poryscriptz.Compile(src, gt)
	if compErr != nil {
		t.Fatalf("[%s] poryscriptZ compile error: %v", tc.name, compErr)
	}
	if len(diags) > 0 {
		var msgs []string
		for _, d := range diags {
			msgs = append(msgs, d.Msg)
		}
		t.Fatalf("[%s] poryscriptZ diagnostics: %s", tc.name, strings.Join(msgs, "; "))
	}
	t.Logf("[%s] compiled .s:\n%s", tc.name, asm)

	// --- write .s to scratch temp dir ---
	tmpDir := t.TempDir()
	sFile := filepath.Join(tmpDir, fmt.Sprintf("cf_%s.s", tc.name))
	oFile := filepath.Join(tmpDir, fmt.Sprintf("cf_%s.o", tc.name))
	if err := os.WriteFile(sFile, []byte(asm), 0644); err != nil {
		t.Fatalf("[%s] write .s: %v", tc.name, err)
	}

	// --- assemble .s → .o with MWAS ---
	// Append -o and input file to the shared base args.
	assembleArgs := make([]string, len(mwasBaseArgs), len(mwasBaseArgs)+3)
	copy(assembleArgs, mwasBaseArgs)
	assembleArgs = append(assembleArgs, "-o", oFile, sFile)

	mwasCmd := exec.Command("wine", append([]string{mwas}, assembleArgs...)...)
	mwasCmd.Dir = hgRoot // run from decomp root so -i . resolves includes
	mwasCmd.Env = append(os.Environ(),
		"WINEPREFIX="+wp,
		"WINEARCH=win32",
	)
	out, runErr := mwasCmd.CombinedOutput()
	if runErr != nil {
		t.Fatalf("[%s] mwasmarm.exe FAILED:\n%s\nerr: %v", tc.name, string(out), runErr)
	}
	if len(out) > 0 {
		t.Logf("[%s] mwasmarm output: %s", tc.name, string(out))
	}
	t.Logf("[%s] PASS: assembles clean", tc.name)
}
