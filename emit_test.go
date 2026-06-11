package poryscriptz

import (
	"os"
	"testing"
)

func TestGoldenStraightLine(t *testing.T) {
	src, _ := os.ReadFile("testdata/straightline.poryz")
	want, _ := os.ReadFile("testdata/straightline.golden")
	got, diags, err := Compile(src, hgssT(t))
	if err != nil || len(diags) != 0 {
		t.Fatalf("compile: err=%v diags=%v", err, diags)
	}
	if got != string(want) {
		t.Fatalf("golden mismatch:\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestGoldenFlagGate(t *testing.T) {
	src, _ := os.ReadFile("testdata/flaggate.poryz")
	want, _ := os.ReadFile("testdata/flaggate.golden")
	got, diags, err := Compile(src, hgssT(t))
	if err != nil || len(diags) != 0 {
		t.Fatalf("compile: err=%v diags=%v", err, diags)
	}
	if got != string(want) {
		t.Fatalf("golden mismatch:\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestGoldenVarGate(t *testing.T) {
	src, _ := os.ReadFile("testdata/vargate.poryz")
	want, _ := os.ReadFile("testdata/vargate.golden")
	got, diags, err := Compile(src, hgssT(t))
	if err != nil || len(diags) != 0 {
		t.Fatalf("compile: err=%v diags=%v", err, diags)
	}
	if got != string(want) {
		t.Fatalf("golden mismatch:\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestGoldenRoute30(t *testing.T) {
	src, _ := os.ReadFile("testdata/route30.poryz")
	want, _ := os.ReadFile("testdata/route30.golden")
	got, diags, err := Compile(src, hgssT(t))
	if err != nil || len(diags) != 0 {
		t.Fatalf("compile: err=%v diags=%v", err, diags)
	}
	if got != string(want) {
		t.Fatalf("golden mismatch:\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestGoldenTwoScripts(t *testing.T) {
	src, _ := os.ReadFile("testdata/twoscripts.poryz")
	want, _ := os.ReadFile("testdata/twoscripts.golden")
	got, diags, err := Compile(src, hgssT(t))
	if err != nil || len(diags) != 0 {
		t.Fatalf("compile: err=%v diags=%v", err, diags)
	}
	if got != string(want) {
		t.Fatalf("golden mismatch:\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestGoldenIfElse(t *testing.T) {
	src, _ := os.ReadFile("testdata/ifelse.poryz")
	want, _ := os.ReadFile("testdata/ifelse.golden")
	got, diags, err := Compile(src, hgssT(t))
	if err != nil || len(diags) != 0 {
		t.Fatalf("compile: err=%v diags=%v", err, diags)
	}
	if got != string(want) {
		t.Fatalf("golden mismatch:\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestGoldenElseIf(t *testing.T) {
	src, _ := os.ReadFile("testdata/elseif.poryz")
	want, _ := os.ReadFile("testdata/elseif.golden")
	got, diags, err := Compile(src, hgssT(t))
	if err != nil || len(diags) != 0 {
		t.Fatalf("compile: err=%v diags=%v", err, diags)
	}
	if got != string(want) {
		t.Fatalf("golden mismatch:\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestGoldenWhile(t *testing.T) {
	src, _ := os.ReadFile("testdata/while.poryz")
	want, _ := os.ReadFile("testdata/while.golden")
	got, diags, err := Compile(src, hgssT(t))
	if err != nil || len(diags) != 0 {
		t.Fatalf("compile: err=%v diags=%v", err, diags)
	}
	if got != string(want) {
		t.Fatalf("golden mismatch:\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestGoldenSwitch(t *testing.T) {
	src, _ := os.ReadFile("testdata/switch.poryz")
	want, _ := os.ReadFile("testdata/switch.golden")
	got, diags, err := Compile(src, hgssT(t))
	if err != nil || len(diags) != 0 {
		t.Fatalf("compile: err=%v diags=%v", err, diags)
	}
	if got != string(want) {
		t.Fatalf("golden mismatch:\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestGoldenMacros(t *testing.T) {
	src, _ := os.ReadFile("testdata/macros.poryz")
	want, _ := os.ReadFile("testdata/macros.golden")
	got, diags, err := Compile(src, hgssT(t))
	if err != nil || len(diags) != 0 {
		t.Fatalf("compile: err=%v diags=%v", err, diags)
	}
	if got != string(want) {
		t.Fatalf("golden mismatch:\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

// TestGoldenWarpLiteralCoords is the regression gate for the downstream bug
// report: "warp cannot express literal tile coords since its var arg type
// expects VAR_* constants".  The bug report was incorrect — poryscriptZ has
// always been arity-only (no arg-kind validation); literal integers and
// arbitrary symbolic constants pass through verbatim.  This test pins that
// behaviour so any future arg-kind enforcement must explicitly opt warp (and
// all scrcmd commands) out before landing.
//
// Fixture: warp(MAP_PALLET_TOWN_OAKS_LAB, 0, 8, 11, DIR_SOUTH)
// Expected output: warp MAP_PALLET_TOWN_OAKS_LAB, 0, 8, 11, DIR_SOUTH
func TestGoldenWarpLiteralCoords(t *testing.T) {
	src, _ := os.ReadFile("testdata/warp_literal_coords.poryz")
	want, _ := os.ReadFile("testdata/warp_literal_coords.golden")
	got, diags, err := Compile(src, hgssT(t))
	if err != nil || len(diags) != 0 {
		t.Fatalf("compile: err=%v diags=%v", err, diags)
	}
	if got != string(want) {
		t.Fatalf("golden mismatch:\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

// TestGoldenGiveitemNoCheck is a targeted golden test for the MacroCall-lowering
// path using giveitem_no_check — the only macro name exclusive to the macro
// table (no scrcmd collision).  A regression in MacroCall lowering would cause
// this fixture to fail distinctly from the scrcmd emit path.
func TestGoldenGiveitemNoCheck(t *testing.T) {
	src, _ := os.ReadFile("testdata/giveitem_no_check.poryz")
	want, _ := os.ReadFile("testdata/giveitem_no_check.golden")
	got, diags, err := Compile(src, hgssT(t))
	if err != nil || len(diags) != 0 {
		t.Fatalf("compile: err=%v diags=%v", err, diags)
	}
	if got != string(want) {
		t.Fatalf("golden mismatch:\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}
