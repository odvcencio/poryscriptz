//go:build decompparity

// Package parity contains a build-tagged COMPARE test that proves a
// poryscriptZ-emitted .s file assembles to bytes identical to the decomp's own
// reference binary for scr_seq_0027.
//
// # Environment requirements
//
//   - Wine prefix: $POKEHG_WINEPREFIX or $WINEPREFIX, defaulting to ~/.wine-pokehg-build
//   - Metrowerks MWAS: $POKEHG_ROOT/tools/mwccarm/2.0/sp2p2/mwasmarm.exe
//   - arm-none-eabi-objcopy in $PATH
//   - $POKEHG_ROOT: explicit path to an isolated pokeheartgold decomp checkout
//
// # Running manually
//
//	go test -tags decompparity ./parity/ -v
//
// # Exact assembler command (from decomp Makefile, run from $POKEHG_ROOT)
//
//	WINEPREFIX=~/.wine-pokehg-build WINEARCH=win32 wine \
//	  tools/mwccarm/2.0/sp2p2/mwasmarm.exe \
//	  -DHEARTGOLD -DGAME_REMASTER=0 -DENGLISH -DPM_KEEP_ASSERTS \
//	  -DSDK_ARM9 -DSDK_CODE_ARM -DSDK_FINALROM -DPM_ASM -DSDK_ASM \
//	  -proc arm5te -g -gccinc \
//	  -i . -i ./include -i ./asm/include -i ./files \
//	  -i ./lib/asm/include -i ./lib/NitroDWC/asm/include \
//	  -i ./lib/MSL_C/asm/include -i ./lib/NitroSDK/asm/include \
//	  -i ./lib/syscall/asm/include -i ./asm -i ./files/msgdata \
//	  -I./lib/include \
//	  -o <tmp>.o <out.s>
//	arm-none-eabi-objcopy -O binary --file-alignment 4 <tmp>.o <tmp>.bin
//	sha1sum <tmp>.bin  # must equal e8cf4832db8e1131f810e3ee00cbece3df9fff08
package parity

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	poryscriptz "github.com/odvcencio/poryscriptz"
	"github.com/odvcencio/poryscriptz/target"
)

// refSHA1 is the sha1 of files/fielddata/script/scr_seq/scr_seq_0027.bin
// from scr_seq.sha1 in the decomp repo.
const refSHA1 = "e8cf4832db8e1131f810e3ee00cbece3df9fff08"

// pokehgRoot requires an explicit isolated decomp checkout.
func pokehgRoot(t *testing.T) string {
	t.Helper()
	if v := os.Getenv("POKEHG_ROOT"); v != "" {
		return v
	}
	t.Skip("decompparity: set POKEHG_ROOT to an isolated decomp checkout")
	return ""
}

// winePrefix returns the Wine prefix path.
func winePrefix() string {
	if v := os.Getenv("POKEHG_WINEPREFIX"); v != "" {
		return v
	}
	if v := os.Getenv("WINEPREFIX"); v != "" {
		return v
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".wine-pokehg-build")
}

// mwasPath lets a separate, read-only decomp checkout use an installed MWAS.
func mwasPath(hgRoot string) string {
	if v := os.Getenv("POKEHG_MWAS"); v != "" {
		return v
	}
	return filepath.Join(hgRoot, "tools", "mwccarm", "2.0", "sp2p2", "mwasmarm.exe")
}

// scrcmdJSON returns the path to the vocab/testdata/scrcmd.json relative to the
// module root (i.e., the parent of this parity package).
func scrcmdJSON(t *testing.T) string {
	t.Helper()
	// __FILE__ is in parity/, parent is module root.
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("decompparity: runtime.Caller failed")
	}
	return filepath.Join(filepath.Dir(thisFile), "..", "vocab", "testdata", "scrcmd.json")
}

// TestCompareScrSeq0027 compiles testdata/parity/scr_seq_0027.poryz via
// poryscriptZ and asserts the assembled binary is byte-identical to the
// decomp's reference sha1 for scr_seq_0027.bin.
//
// Chosen script: scr_seq_0027 — a single-script file with one "end" command.
// It is one of the smallest real scripts (13 lines), uses no message ops,
// no movement tables, and no constructs beyond what poryscriptZ v0.1 supports.
//
// poryscriptZ expressibility of this script: FULL.
// The original includes event_0027.h (defines _EV_scr_seq_0027_000 = 0) and
// msg_0003_EVERYWHERE.h. Neither define is referenced in the script body, so
// they are noise-free for assembly and the poryscriptZ output (without those
// includes) assembles to the same bytes.
func TestCompareScrSeq0027(t *testing.T) {
	// --- prerequisite checks ---
	hgRoot := pokehgRoot(t)
	mwas := mwasPath(hgRoot)
	winePrefix := winePrefix()

	for _, check := range []struct {
		label string
		path  string
	}{
		{"decomp root", hgRoot},
		{"mwasmarm.exe", mwas},
		{"wine prefix", winePrefix},
	} {
		if _, err := os.Stat(check.path); err != nil {
			t.Skipf("decompparity: prerequisite not found (%s=%s): %v", check.label, check.path, err)
		}
	}
	if _, err := exec.LookPath("wine"); err != nil {
		t.Skip("decompparity: wine not in PATH")
	}
	if _, err := exec.LookPath("arm-none-eabi-objcopy"); err != nil {
		t.Skip("decompparity: arm-none-eabi-objcopy not in PATH")
	}

	// --- compile .poryz → .s ---
	_, thisFile, _, _ := runtime.Caller(0)
	parityTestdata := filepath.Join(filepath.Dir(thisFile), "..", "testdata", "parity", "scr_seq_0027.poryz")
	src, err := os.ReadFile(parityTestdata)
	if err != nil {
		t.Fatalf("read .poryz: %v", err)
	}

	scrcmd := scrcmdJSON(t)
	tgt, err := target.HGSS(scrcmd)
	if err != nil {
		t.Fatalf("load HGSS target: %v", err)
	}

	asm, diags, err := poryscriptz.Compile(src, tgt)
	if err != nil {
		t.Fatalf("poryscriptZ compile: %v", err)
	}
	if len(diags) > 0 {
		var msgs []string
		for _, d := range diags {
			msgs = append(msgs, d.Msg)
		}
		t.Fatalf("poryscriptZ diagnostics: %s", strings.Join(msgs, "; "))
	}

	// --- write .s to temp dir ---
	tmpDir := t.TempDir()
	sFile := filepath.Join(tmpDir, "scr_seq_0027_parity.s")
	if err := os.WriteFile(sFile, []byte(asm), 0644); err != nil {
		t.Fatalf("write .s: %v", err)
	}
	t.Logf("poryscriptZ .s output:\n%s", asm)

	// --- assemble .s → .o ---
	// Command mirrors the decomp Makefile exactly (see scr_seq.mk + common.mk):
	//   $(WINE) $(MWAS) $(MWASFLAGS) -DPM_ASM -o $*.o $<
	// MWASFLAGS = -DHEARTGOLD -DGAME_REMASTER=0 -DENGLISH -DPM_KEEP_ASSERTS
	//             -DSDK_ARM9 -DSDK_CODE_ARM -DSDK_FINALROM
	//             -proc arm5te -g -gccinc
	//             -i . -i ./include -i ./asm/include -i ./files
	//             -i ./lib/asm/include -i ./lib/NitroDWC/asm/include
	//             -i ./lib/MSL_C/asm/include -i ./lib/NitroSDK/asm/include
	//             -i ./lib/syscall/asm/include -i ./asm -i ./files/msgdata
	//             -I./lib/include -DSDK_ASM
	oFile := filepath.Join(tmpDir, "scr_seq_0027_parity.o")
	assembleArgs := []string{
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
		"-o", oFile,
		sFile,
	}

	mwasCmd := exec.Command("wine", append([]string{mwas}, assembleArgs...)...)
	mwasCmd.Dir = hgRoot // must run from decomp root so -i . resolves headers
	mwasCmd.Env = append(os.Environ(),
		"WINEPREFIX="+winePrefix,
		"WINEARCH=win32",
	)
	out, err := mwasCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("mwasmarm.exe failed:\n%s\nerr: %v", string(out), err)
	}
	if len(out) > 0 {
		t.Logf("mwasmarm output: %s", string(out))
	}

	// --- objcopy .o → .bin ---
	binFile := filepath.Join(tmpDir, "scr_seq_0027_parity.bin")
	objcopyCmd := exec.Command("arm-none-eabi-objcopy",
		"-O", "binary", "--file-alignment", "4",
		oFile, binFile,
	)
	if out, err := objcopyCmd.CombinedOutput(); err != nil {
		t.Fatalf("arm-none-eabi-objcopy failed:\n%s\nerr: %v", string(out), err)
	}

	// --- sha1 and compare ---
	got, err := fileSHA1(binFile)
	if err != nil {
		t.Fatalf("sha1 of output bin: %v", err)
	}
	t.Logf("assembled binary sha1: %s", got)
	t.Logf("reference sha1:        %s", refSHA1)

	if got != refSHA1 {
		t.Fatalf("BYTE MISMATCH: poryscriptZ output assembled to %s, want %s", got, refSHA1)
	}
	t.Logf("PASS: byte-identical to decomp reference binary")
}

func fileSHA1(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
