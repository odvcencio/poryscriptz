package poryscriptz

import "fmt"

// Instr is a single flat instruction produced by Lower.
//
// If Label is non-empty, the emitter renders "<Label>:\n" before the body.
// If Opcode is -1 the Instr is a label-only sentinel (script entry point or
// branch target) and the emitter skips the macro body line.
//
// Macro, when non-empty, overrides Opcode→MacroName resolution in Emit and is
// used as the literal macro name.  This is the escape hatch for control-flow
// wrapper macros (goto_if_unset, goto_if_ne, compare_var_to_value) that are
// not base opcodes in scrcmd.json and therefore have no Opcode index.
// TODO(task-N): route these through GameTarget when a second game is added.
type Instr struct {
	Label  string   // entry-point or branch-target label; non-empty on sentinels
	Opcode int      // -1 for label-only sentinels; otherwise the vocab opcode
	Macro  string   // literal macro override; takes precedence over Opcode in Emit
	Args   []string // raw arg text (empty slice for zero-arg commands)
}

// lowerCtx carries per-script mutable state for Lower.
// labelSeq is a pointer so it can be shared across all scripts in a file,
// making branch labels file-level monotonic and preventing duplicate labels
// when multiple branching scripts are compiled together.
type lowerCtx struct {
	out      []Instr
	labelSeq *int // file-level monotonic counter; yields _L0, _L1, … across all scripts
}

func (c *lowerCtx) nextLabel() string {
	l := fmt.Sprintf("_L%d", *c.labelSeq)
	*c.labelSeq++
	return l
}

func (c *lowerCtx) append(i Instr) { c.out = append(c.out, i) }

// lowerStmts recursively lowers a slice of Stmts into c.out.
func (c *lowerCtx) lowerStmts(stmts []Stmt) {
	for _, stmt := range stmts {
		switch st := stmt.(type) {
		case Call:
			c.append(Instr{
				Opcode: st.Cmd.Opcode,
				Args:   st.Args,
			})
		case If:
			c.lowerIf(st)
		}
	}
}

// lowerIf lowers an If statement into a branch + body + label-sentinel pattern.
//
// flag:  goto_if_unset FLAG, _Lx  ; <body> ; _Lx:
// vareq: compare_var_to_value VAR, N ; goto_if_ne _Lx ; <body> ; _Lx:
//
// The wrapper macro names (goto_if_unset, goto_if_ne, compare_var_to_value) are
// HGSS-specific and emitted as literal Macro strings.
// TODO(task-N): route through GameTarget when a second game target is added.
func (c *lowerCtx) lowerIf(st If) {
	label := c.nextLabel()
	switch st.Kind {
	case "flag":
		c.append(Instr{Macro: "goto_if_unset", Args: []string{st.Flag, label}})
	case "vareq":
		c.append(Instr{Macro: "compare_var_to_value", Args: []string{st.Var, st.Value}})
		c.append(Instr{Macro: "goto_if_ne", Args: []string{label}})
	}
	c.lowerStmts(st.Body)
	// Branch target sentinel (label only, no macro body).
	c.append(Instr{Label: label, Opcode: -1})
}

// Lower converts a resolved Program into a flat Instr slice.
//
// Layout per script:
//  1. A label sentinel  Instr{Label: script.Name, Opcode: -1}
//  2. Instrs for each statement (Call or If)
//
// The Emit function translates Opcode → macro string via tgt.MacroName so that
// unnamed opcodes render as "scrcmd_NNN" without Lower knowing about the target.
// Instrs with a non-empty Macro field bypass that lookup and use Macro directly.
func Lower(p *Program) []Instr {
	var out []Instr
	// seq is shared across all scripts so that branch labels are file-level
	// monotonic (_L0, _L1, …) and never collide between scripts in one file.
	seq := 0
	for _, s := range p.Scripts {
		ctx := &lowerCtx{labelSeq: &seq}
		// Entry-point sentinel — label only, no macro body.
		ctx.append(Instr{Label: s.Name, Opcode: -1})
		ctx.lowerStmts(s.Stmts)
		out = append(out, ctx.out...)
	}
	return out
}
