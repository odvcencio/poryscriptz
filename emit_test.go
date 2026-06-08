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
