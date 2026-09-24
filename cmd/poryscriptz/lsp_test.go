package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func lspFrame(v any) string {
	b, _ := json.Marshal(v)
	return fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(b), b)
}
func TestLSPDiagnosticsNavigationAndHighlight(t *testing.T) {
	source := "package m\nscript Start {\n    GoTo(_Done)\n}\nlabel _Done {\n    End()\n}\n"
	uri := "file:///tmp/example.poryz"
	messages := []any{
		map[string]any{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": map[string]any{}},
		map[string]any{"jsonrpc": "2.0", "method": "textDocument/didOpen", "params": map[string]any{"textDocument": map[string]any{"uri": uri, "text": source, "version": 1}}},
		map[string]any{"jsonrpc": "2.0", "id": 2, "method": "textDocument/hover", "params": map[string]any{"textDocument": map[string]any{"uri": uri}, "position": map[string]any{"line": 2, "character": 6}}},
		map[string]any{"jsonrpc": "2.0", "id": 3, "method": "textDocument/definition", "params": map[string]any{"textDocument": map[string]any{"uri": uri}, "position": map[string]any{"line": 2, "character": 11}}},
		map[string]any{"jsonrpc": "2.0", "id": 4, "method": "textDocument/semanticTokens/full", "params": map[string]any{"textDocument": map[string]any{"uri": uri}}},
		map[string]any{"jsonrpc": "2.0", "id": 5, "method": "textDocument/completion", "params": map[string]any{"textDocument": map[string]any{"uri": uri}, "position": map[string]any{"line": 2, "character": 11}}},
	}
	var input strings.Builder
	for _, m := range messages {
		input.WriteString(lspFrame(m))
	}
	var output, errors bytes.Buffer
	if code := runLSP(strings.NewReader(input.String()), &output, &errors); code != 0 {
		t.Fatalf("LSP exit %d: %s", code, errors.String())
	}
	got := output.String()
	for _, want := range []string{`"semanticTokensProvider"`, `"diagnostics":[]`, `GoTo(dest)`, `"line":4`, `"data":[`} {
		if !strings.Contains(got, want) {
			t.Errorf("response missing %q", want)
		}
	}
	items:=completion("ITEM_PO").([]any)
	if len(items)==0 {t.Fatal("constant completion returned no items")}
}
