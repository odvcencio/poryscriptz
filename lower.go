package poryscriptz

// Instr is a single flat instruction produced by Lower.
//
// If Label is non-empty, the emitter renders "<Label>:\n" before the body.
// If Opcode is -1 the Instr is a label-only sentinel (script entry point) and
// the emitter skips the macro body line.
//
// Task 6 will extend Lower to emit branch/label Instrs for if-statements.
type Instr struct {
	Label  string   // entry-point label, e.g. "Greeter"; non-empty only on entry sentinels
	Opcode int      // -1 for label-only sentinels; otherwise the vocab opcode
	Args   []string // raw arg text (empty slice for zero-arg commands)
}

// Lower converts a resolved Program into a flat Instr slice.
//
// Layout per script:
//   1. A label sentinel  Instr{Label: script.Name, Opcode: -1}
//   2. One Instr per Call statement
//
// The Emit function translates Opcode → macro string via tgt.MacroName so that
// unnamed opcodes render as "scrcmd_NNN" without Lower knowing about the target.
func Lower(p *Program) []Instr {
	var out []Instr
	for _, s := range p.Scripts {
		// Entry-point sentinel — label only, no macro body.
		out = append(out, Instr{Label: s.Name, Opcode: -1})
		for _, stmt := range s.Stmts {
			switch st := stmt.(type) {
			case Call:
				out = append(out, Instr{
					Opcode: st.Cmd.Opcode,
					Args:   st.Args,
				})
			}
		}
	}
	return out
}
