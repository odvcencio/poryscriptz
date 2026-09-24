package poryscriptz

import "github.com/odvcencio/poryscriptz/target"

// Compile parses, resolves, lowers, and emits src for the given target.
// Returns the assembled text, any diagnostics (on parse/resolve error), and a
// hard parse error if the source is structurally unparseable.
func Compile(src []byte, tgt target.GameTarget) (string, []Diag, error) {
	root, w, err := Parse(src)
	if err != nil {
		return "", nil, err
	}
	prog, diags := Resolve(w, root, tgt)
	if len(diags) != 0 {
		return "", diags, nil
	}
	return Emit(prog, tgt), nil, nil
}
