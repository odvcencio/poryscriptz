package poryscriptz

import "m31labs.dev/poryscriptz/target"

// Compile (Task 4: parse + resolve only). Task 5 extends this to lower + emit.
func Compile(src []byte, tgt target.GameTarget) (string, []Diag, error) {
	root, w, err := Parse(src)
	if err != nil {
		return "", nil, err
	}
	prog, diags := Resolve(w, root, tgt)
	_ = prog
	return "", diags, nil // TODO(Task 5): if no diags, Lower(prog) then Emit(...)
}
