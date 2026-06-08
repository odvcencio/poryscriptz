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
