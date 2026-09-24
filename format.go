package poryscriptz

import (
	"bytes"
	"strings"
)

// Format applies a stable four-space layout to valid poryz source.
func Format(src []byte) ([]byte, error) {
	if _, _, err := Parse(src); err != nil {
		return nil, err
	}
	var out strings.Builder
	depth, blanks := 0, 0
	for _, raw := range strings.Split(strings.ReplaceAll(string(src), "\r\n", "\n"), "\n") {
		line := strings.TrimSpace(raw)
		if line == "" {
			if out.Len() > 0 {
				blanks++
			}
			continue
		}
		if blanks > 0 && out.Len() > 0 {
			out.WriteByte('\n')
		}
		blanks = 0
		closed := 0
		for closed < len(line) && line[closed] == '}' {
			closed++
		}
		lineDepth := depth - closed
		if lineDepth < 0 {
			lineDepth = 0
		}
		isCase := strings.HasPrefix(line, "case ") || strings.HasPrefix(line, "default:")
		indent := lineDepth
		if isCase && indent > 0 {
			indent--
		}
		out.WriteString(strings.Repeat("    ", indent))
		out.WriteString(line)
		out.WriteByte('\n')
		open, close := countBraces(line)
		depth += open - close
		if depth < 0 {
			depth = 0
		}
	}
	return bytes.TrimRight([]byte(out.String()), "\n"), nil
}

func countBraces(line string) (open, close int) {
	quoted, escaped := false, false
	for i := 0; i < len(line); i++ {
		c := line[i]
		if !quoted && i+1 < len(line) && c == '/' && line[i+1] == '/' {
			break
		}
		if c == '"' && !escaped {
			quoted = !quoted
		}
		if !quoted {
			if c == '{' {
				open++
			}
			if c == '}' {
				close++
			}
		}
		if c == '\\' && !escaped {
			escaped = true
		} else {
			escaped = false
		}
	}
	return
}
