package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode/utf16"

	gts "github.com/odvcencio/gotreesitter"
	poryz "github.com/odvcencio/poryscriptz"
	"github.com/odvcencio/poryscriptz/target"
	"github.com/odvcencio/poryscriptz/vocab"
)

type rpcMessage struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
}
type lspPosition struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}
type lspRange struct {
	Start lspPosition `json:"start"`
	End   lspPosition `json:"end"`
}
type lspDoc struct {
	URI     string `json:"uri"`
	Text    string `json:"text"`
	Version int    `json:"version"`
}
type lspServer struct {
	in     *bufio.Reader
	out    io.Writer
	errors io.Writer
	docs   map[string]lspDoc
	target target.GameTarget
}

func runLSP(input io.Reader, output io.Writer, errors io.Writer) int {
	tgt, err := target.HGSS("")
	if err != nil {
		fmt.Fprintln(errors, err)
		return 1
	}
	s := &lspServer{in: bufio.NewReader(input), out: output, errors: errors, docs: map[string]lspDoc{}, target: tgt}
	for {
		payload, err := s.read()
		if err == io.EOF {
			return 0
		}
		if err != nil {
			fmt.Fprintln(errors, "poryz lsp:", err)
			return 1
		}
		var msg rpcMessage
		if err := json.Unmarshal(payload, &msg); err != nil {
			fmt.Fprintln(errors, "poryz lsp: invalid JSON:", err)
			continue
		}
		if msg.Method == "exit" {
			return 0
		}
		s.handle(msg)
	}
}

func (s *lspServer) read() ([]byte, error) {
	length := -1
	for {
		line, err := s.in.ReadString('\n')
		if err != nil {
			return nil, err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		key, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		if strings.EqualFold(key, "Content-Length") {
			length, err = strconv.Atoi(strings.TrimSpace(value))
			if err != nil {
				return nil, err
			}
		}
	}
	if length < 0 || length > 16*1024*1024 {
		return nil, fmt.Errorf("invalid Content-Length %d", length)
	}
	data := make([]byte, length)
	_, err := io.ReadFull(s.in, data)
	return data, err
}
func (s *lspServer) send(value any) {
	data, err := json.Marshal(value)
	if err != nil {
		fmt.Fprintln(s.errors, err)
		return
	}
	fmt.Fprintf(s.out, "Content-Length: %d\r\n\r\n", len(data))
	_, _ = s.out.Write(data)
}
func (s *lspServer) respond(id json.RawMessage, result any) {
	if len(id) > 0 {
		s.send(map[string]any{"jsonrpc": "2.0", "id": id, "result": result})
	}
}

func (s *lspServer) handle(msg rpcMessage) {
	switch msg.Method {
	case "initialize":
		s.respond(msg.ID, map[string]any{"capabilities": map[string]any{
			"textDocumentSync": 1, "hoverProvider": true, "definitionProvider": true,
			"completionProvider":     map[string]any{"triggerCharacters": []string{"_"}},
			"semanticTokensProvider": map[string]any{"full": true, "legend": map[string]any{"tokenTypes": []string{"keyword", "comment", "string", "number", "function", "type", "variable", "enumMember"}, "tokenModifiers": []string{}}},
		}, "serverInfo": map[string]string{"name": "poryz", "version": "0.3.0"}})
	case "initialized":
	case "shutdown":
		s.respond(msg.ID, nil)
	case "textDocument/didOpen":
		var p struct {
			TextDocument lspDoc `json:"textDocument"`
		}
		if json.Unmarshal(msg.Params, &p) != nil {
			return
		}
		s.docs[p.TextDocument.URI] = p.TextDocument
		s.publish(p.TextDocument)
	case "textDocument/didChange":
		var p struct {
			TextDocument struct {
				URI     string `json:"uri"`
				Version int    `json:"version"`
			} `json:"textDocument"`
			ContentChanges []struct {
				Text string `json:"text"`
			} `json:"contentChanges"`
		}
		if json.Unmarshal(msg.Params, &p) != nil || len(p.ContentChanges) == 0 {
			return
		}
		doc := lspDoc{URI: p.TextDocument.URI, Version: p.TextDocument.Version, Text: p.ContentChanges[len(p.ContentChanges)-1].Text}
		s.docs[doc.URI] = doc
		s.publish(doc)
	case "textDocument/didClose":
		var p struct {
			TextDocument struct {
				URI string `json:"uri"`
			} `json:"textDocument"`
		}
		if json.Unmarshal(msg.Params, &p) != nil {
			return
		}
		delete(s.docs, p.TextDocument.URI)
		s.send(map[string]any{"jsonrpc": "2.0", "method": "textDocument/publishDiagnostics", "params": map[string]any{"uri": p.TextDocument.URI, "diagnostics": []any{}}})
	case "textDocument/hover", "textDocument/definition", "textDocument/completion":
		var p struct {
			TextDocument struct {
				URI string `json:"uri"`
			} `json:"textDocument"`
			Position lspPosition `json:"position"`
		}
		if json.Unmarshal(msg.Params, &p) != nil {
			s.respond(msg.ID, nil)
			return
		}
		doc, ok := s.docs[p.TextDocument.URI]
		if !ok {
			s.respond(msg.ID, nil)
			return
		}
		word := wordAt(doc.Text, p.Position)
		switch msg.Method {
		case "textDocument/hover":
			s.respond(msg.ID, s.hover(word))
		case "textDocument/definition":
			s.respond(msg.ID, definition(doc, word))
		case "textDocument/completion":
			s.respond(msg.ID, completion(word))
		}
	case "textDocument/semanticTokens/full":
		var p struct {
			TextDocument struct {
				URI string `json:"uri"`
			} `json:"textDocument"`
		}
		if json.Unmarshal(msg.Params, &p) != nil {
			s.respond(msg.ID, map[string]any{"data": []int{}})
			return
		}
		doc, ok := s.docs[p.TextDocument.URI]
		if !ok {
			s.respond(msg.ID, map[string]any{"data": []int{}})
			return
		}
		s.respond(msg.ID, map[string]any{"data": semanticTokens(doc.Text)})
	default:
		if len(msg.ID) > 0 {
			s.send(map[string]any{"jsonrpc": "2.0", "id": msg.ID, "error": map[string]any{"code": -32601, "message": "method not supported"}})
		}
	}
}

func semanticTokens(src string) []int {
	ranges, err := poryz.Highlight([]byte(src))
	if err != nil {
		return []int{}
	}
	types := map[string]int{"keyword": 0, "comment": 1, "string": 2, "number": 3, "function": 4, "type": 5, "label": 6, "constant": 7}
	type token struct{ line, col, length, kind int }
	var tokens []token
	for _, r := range ranges {
		kind, ok := types[r.Capture]
		if !ok {
			continue
		}
		start, end := int(r.StartByte), int(r.EndByte)
		for start < end {
			stop := end
			if i := strings.IndexByte(src[start:end], '\n'); i >= 0 {
				stop = start + i
			}
			if stop > start {
				line, col := bytePosition(src, start)
				length := len(utf16.Encode([]rune(src[start:stop])))
				tokens = append(tokens, token{line, col, length, kind})
			}
			if stop == end {
				break
			}
			start = stop + 1
		}
	}
	sort.Slice(tokens, func(i, j int) bool {
		if tokens[i].line != tokens[j].line {
			return tokens[i].line < tokens[j].line
		}
		return tokens[i].col < tokens[j].col
	})
	data := make([]int, 0, len(tokens)*5)
	prevLine, prevCol := 0, 0
	for _, t := range tokens {
		dl := t.line - prevLine
		dc := t.col
		if dl == 0 {
			dc -= prevCol
		}
		data = append(data, dl, dc, t.length, t.kind, 0)
		prevLine, prevCol = t.line, t.col
	}
	return data
}
func bytePosition(src string, offset int) (int, int) {
	if offset > len(src) {
		offset = len(src)
	}
	before := src[:offset]
	line := strings.Count(before, "\n")
	start := strings.LastIndexByte(before, '\n') + 1
	return line, len(utf16.Encode([]rune(before[start:])))
}

func (s *lspServer) publish(doc lspDoc) {
	_, diags, err := poryz.Compile([]byte(doc.Text), s.target)
	items := make([]any, 0)
	if err != nil {
		line, col, msg := parseSyntaxError(err)
		items = append(items, diagnostic(doc.Text, line, col, msg))
	} else {
		for _, d := range diags {
			items = append(items, diagnostic(doc.Text, d.Line, d.Col, d.Msg))
		}
	}
	s.send(map[string]any{"jsonrpc": "2.0", "method": "textDocument/publishDiagnostics", "params": map[string]any{"uri": doc.URI, "version": doc.Version, "diagnostics": items}})
}

var syntaxPos = regexp.MustCompile(`^(\d+):(\d+):\s*(.*)$`)

func parseSyntaxError(err error) (int, int, string) {
	m := syntaxPos.FindStringSubmatch(err.Error())
	if m == nil {
		return 1, 1, err.Error() + "; check braces and commas"
	}
	line, _ := strconv.Atoi(m[1])
	col, _ := strconv.Atoi(m[2])
	return line, col, m[3] + "; check braces and commas"
}
func diagnostic(src string, line, col int, msg string) any {
	start := lspPosition{Line: max(0, line-1), Character: byteColToUTF16(src, line, col)}
	return map[string]any{"range": lspRange{Start: start, End: lspPosition{Line: start.Line, Character: start.Character + 1}}, "severity": 1, "source": "poryz", "message": msg}
}
func byteColToUTF16(src string, line, col int) int {
	lines := strings.Split(src, "\n")
	if line < 1 || line > len(lines) {
		return max(0, col-1)
	}
	b := []byte(lines[line-1])
	n := min(max(0, col-1), len(b))
	return len(utf16.Encode([]rune(string(b[:n]))))
}
func positionByte(src string, p lspPosition) int {
	lines := strings.SplitAfter(src, "\n")
	if p.Line < 0 || p.Line >= len(lines) {
		return len(src)
	}
	offset := 0
	for i := 0; i < p.Line; i++ {
		offset += len(lines[i])
	}
	line := lines[p.Line]
	units := 0
	for i, r := range line {
		if units >= p.Character {
			return offset + i
		}
		units += len(utf16.Encode([]rune{r}))
	}
	return offset + len(line)
}
func wordAt(src string, p lspPosition) string {
	offset := positionByte(src, p)
	if offset > len(src) {
		offset = len(src)
	}
	start, end := offset, offset
	for start > 0 && wordByte(src[start-1]) {
		start--
	}
	for end < len(src) && wordByte(src[end]) {
		end++
	}
	return src[start:end]
}
func wordByte(b byte) bool {
	return b == '_' || b >= 'A' && b <= 'Z' || b >= 'a' && b <= 'z' || b >= '0' && b <= '9'
}

func (s *lspServer) hover(word string) any {
	if word == "" {
		return nil
	}
	if m, ok := s.target.Macros().ByName(word); ok {
		params := m.Params
		if len(params) == 0 && len(m.Args) > 0 {
			params = make([]string, len(m.Args))
			for i := range params {
				params[i] = "arg"
			}
		}
		kind := "Script command"
		if m.Movement {
			kind = "Movement command"
		}
		doc := kind + " from the HGSS macro vocabulary."
		if m.Doc != "" {
			doc = m.Doc
		}
		return map[string]any{"contents": map[string]string{"kind": "markdown", "value": fmt.Sprintf("```poryz\n%s(%s)\n```\n%s", m.Name, strings.Join(params, ", "), doc)}}
	}
	if c, ok := s.target.Vocabulary().ByName(word); ok {
		return map[string]any{"contents": map[string]string{"kind": "markdown", "value": fmt.Sprintf("`%s`: script opcode %d, %d arguments", c.Name, c.Opcode, len(c.Args))}}
	}
	return nil
}
func definition(doc lspDoc, word string) any {
	if word == "" {
		return nil
	}
	root, w, err := poryz.Parse([]byte(doc.Text))
	if err != nil {
		return nil
	}
	var found any
	var search func(*gts.Node)
	search = func(n *gts.Node) {
		if n == nil || found != nil {
			return
		}
		switch w.Type(n) {
		case "script_declaration", "label_declaration", "movement_declaration":
			name := w.Field(n, "name")
			if w.Text(name) == word {
				line, col := w.Pos(name)
				start := lspPosition{Line: line - 1, Character: byteColToUTF16(doc.Text, line, col)}
				found = map[string]any{"uri": doc.URI, "range": lspRange{Start: start, End: lspPosition{Line: start.Line, Character: start.Character + len(utf16.Encode([]rune(word)))}}}
				return
			}
		}
		for i := 0; i < n.ChildCount(); i++ {
			search(n.Child(i))
		}
	}
	search(root)
	return found
}
func completion(prefix string) any {
	if prefix == "" {
		return []any{}
	}
	all, err := vocab.Constants(prefix)
	if err != nil {
		return []any{}
	}
	items := make([]any, 0, min(len(all), 200))
	for _, c := range all {
		if len(items) == 200 {
			break
		}
		items = append(items, map[string]any{"label": c.Name, "kind": 21, "detail": c.Kind + " = " + c.Value})
	}
	return items
}
