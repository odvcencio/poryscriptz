package poryscriptz_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	poryz "github.com/odvcencio/poryscriptz"
	"github.com/odvcencio/poryscriptz/target"
)

func TestPublishedExamplesCompile(t *testing.T) {
	tgt, err := target.HGSS("")
	if err != nil {
		t.Fatal(err)
	}
	files, err := filepath.Glob("examples/*.poryz")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) < 20 {
		t.Fatalf("want at least 20 recipes, got %d", len(files))
	}
	for _, file := range files {
		t.Run(filepath.Base(file), func(t *testing.T) {
			src, err := os.ReadFile(file)
			if err != nil {
				t.Fatal(err)
			}
			assertCompiles(t, src, tgt)
		})
	}
	docs := []string{"README.md", "docs/first-script.md", "docs/cookbook.md", "docs/patterns.md"}
	for _, file := range docs {
		data, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		parts := strings.Split(string(data), "```poryz\n")
		for i, part := range parts[1:] {
			src, _, ok := strings.Cut(part, "\n```")
			if !ok {
				t.Fatalf("%s block %d not closed", file, i+1)
			}
			t.Run(filepath.Base(file)+"/block"+string(rune('A'+i)), func(t *testing.T) { assertCompiles(t, []byte(src), tgt) })
		}
	}
}
func assertCompiles(t *testing.T, src []byte, tgt target.GameTarget) {
	t.Helper()
	_, diags, err := poryz.Compile(src, tgt)
	if err != nil {
		t.Fatal(err)
	}
	if len(diags) > 0 {
		t.Fatalf("%+v", diags)
	}
}
