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
	Label     string   // entry-point or branch-target label; non-empty on sentinels
	Opcode    int      // -1 for label-only sentinels; otherwise the vocab opcode
	Macro     string   // literal macro override; takes precedence over Opcode in Emit
	Args      []string // raw arg text (empty slice for zero-arg commands)
	Directive string   // assembler directive, emitted without a command lookup
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
		case MacroCall:
			name := st.Macro.Name
			if st.Macro.EmitName != "" {
				name = st.Macro.EmitName
			}
			// Macro calls emit verbatim: "name arg1, arg2, …"
			// The decomp assembler expands them via asm/macros/script.inc.
			c.append(Instr{
				Opcode: -1,
				Macro:  name,
				Args:   st.Args,
			})
		case PatternCall:
			for _, ins := range expandPattern(st.Name, st.Args) {
				c.append(ins)
			}
		case If:
			c.lowerIf(st)
		case While:
			c.lowerWhile(st)
		case Switch:
			c.lowerSwitch(st)
		}
	}
}

// lowerIf lowers an If statement into a branch + body + label-sentinel pattern.
//
// Without else (v0.1 behaviour, unchanged):
//
//	flag:  goto_if_unset FLAG, _Lx  ; <body> ; _Lx:
//	vareq: compare_var_to_value VAR, N ; goto_if_ne _Lx ; <body> ; _Lx:
//
// With else:
//
//	flag:  goto_if_unset FLAG, _Lelse ; <then> ; goto _Lend ; _Lelse: ; <else> ; _Lend:
//	vareq: compare_var_to_value VAR,N ; goto_if_ne _Lelse ; <then> ; goto _Lend ; _Lelse: ; <else> ; _Lend:
//
// Each call allocates fresh labels from the file-monotonic counter so that
// nested and sibling constructs never produce duplicate label names.
//
// The wrapper macro names (goto_if_unset, goto_if_ne, compare_var_to_value) are
// HGSS-specific and emitted as literal Macro strings.
// TODO(task-N): route through GameTarget when a second game target is added.
func (c *lowerCtx) lowerIf(st If) {
	if len(st.Else) == 0 {
		// No else — v0.1 single-skip-label behaviour.
		skipLabel := c.nextLabel()
		switch st.Kind {
		case "flag":
			c.append(Instr{Macro: "goto_if_unset", Args: []string{st.Flag, skipLabel}})
		case "vareq":
			c.append(Instr{Macro: "compare_var_to_value", Args: []string{st.Var, st.Value}})
			c.append(Instr{Macro: "goto_if_ne", Args: []string{skipLabel}})
		}
		c.lowerStmts(st.Body)
		c.append(Instr{Label: skipLabel, Opcode: -1})
		return
	}

	// Has else — allocate _Lelse (jump target when condition is false) and
	// _Lend (jump target after the then-body to skip the else-body).
	lelse := c.nextLabel()
	lend := c.nextLabel()

	switch st.Kind {
	case "flag":
		c.append(Instr{Macro: "goto_if_unset", Args: []string{st.Flag, lelse}})
	case "vareq":
		c.append(Instr{Macro: "compare_var_to_value", Args: []string{st.Var, st.Value}})
		c.append(Instr{Macro: "goto_if_ne", Args: []string{lelse}})
	}
	c.lowerStmts(st.Body)
	// Unconditional jump past the else-body.
	c.append(Instr{Macro: "goto", Args: []string{lend}})
	// Else label sentinel + else body.
	c.append(Instr{Label: lelse, Opcode: -1})
	c.lowerStmts(st.Else)
	// End label sentinel.
	c.append(Instr{Label: lend, Opcode: -1})
}

// lowerWhile lowers a While statement into a top-test loop:
//
//	_Ltop:
//	  <¬cond jump → _Lend>
//	  <body>
//	  goto _Ltop
//	_Lend:
//
// The condition is negated in the same way as lowerIf's no-else path:
//
//	flag:  goto_if_unset FLAG, _Lend
//	vareq: compare_var_to_value VAR, N ; goto_if_ne _Lend
func (c *lowerCtx) lowerWhile(st While) {
	ltop := c.nextLabel()
	lend := c.nextLabel()

	// Top-of-loop label sentinel.
	c.append(Instr{Label: ltop, Opcode: -1})

	// Negated condition jump to lend.
	switch st.Kind {
	case "flag":
		c.append(Instr{Macro: "goto_if_unset", Args: []string{st.Flag, lend}})
	case "vareq":
		c.append(Instr{Macro: "compare_var_to_value", Args: []string{st.Var, st.Value}})
		c.append(Instr{Macro: "goto_if_ne", Args: []string{lend}})
	}

	// Loop body.
	c.lowerStmts(st.Body)

	// Back-edge jump.
	c.append(Instr{Macro: "goto", Args: []string{ltop}})

	// End-of-loop label sentinel.
	c.append(Instr{Label: lend, Opcode: -1})
}

// lowerSwitch lowers a Switch statement into a compare-chain dispatch:
//
//	compare_var_to_value VAR, case0.Val
//	goto_if_eq _L<case0>
//	compare_var_to_value VAR, case1.Val
//	goto_if_eq _L<case1>
//	…
//	goto _Ldefault        (or goto _Lend if no default)
//	_L<case0>:
//	  <body0>
//	  goto _Lend
//	_L<case1>:
//	  <body1>
//	  goto _Lend
//	…
//	_Ldefault:
//	  <default body>
//	_Lend:
//
// Arms do NOT fall through; each arm ends with an unconditional goto _Lend.
// The default arm (if present) falls through naturally into _Lend.
func (c *lowerCtx) lowerSwitch(st Switch) {
	// Allocate one label per case arm, one for default, one for end.
	caseLabels := make([]string, len(st.Cases))
	for i := range st.Cases {
		caseLabels[i] = c.nextLabel()
	}
	var ldefault string
	if len(st.Default) > 0 {
		ldefault = c.nextLabel()
	}
	lend := c.nextLabel()

	// Emit compare-chain.
	for i, sc := range st.Cases {
		c.append(Instr{Macro: "compare_var_to_value", Args: []string{st.Var, sc.Val}})
		c.append(Instr{Macro: "goto_if_eq", Args: []string{caseLabels[i]}})
	}
	// Unconditional jump to default or end.
	if ldefault != "" {
		c.append(Instr{Macro: "goto", Args: []string{ldefault}})
	} else {
		c.append(Instr{Macro: "goto", Args: []string{lend}})
	}

	// Emit case arm bodies.
	for i, sc := range st.Cases {
		c.append(Instr{Label: caseLabels[i], Opcode: -1})
		c.lowerStmts(sc.Body)
		c.append(Instr{Macro: "goto", Args: []string{lend}})
	}

	// Emit default body (falls through to lend).
	if ldefault != "" {
		c.append(Instr{Label: ldefault, Opcode: -1})
		c.lowerStmts(st.Default)
	}

	// End-of-switch label sentinel.
	c.append(Instr{Label: lend, Opcode: -1})
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
	sections := p.Sections
	if len(sections) == 0 {
		for _, s := range p.Scripts {
			sections = append(sections, Section{Name: s.Name, Stmts: s.Stmts})
		}
	}
	for _, s := range sections {
		if s.SourceInclude != "" {
			out = append(out, Instr{Opcode: -1, Directive: fmt.Sprintf(".include %q", s.SourceInclude)})
			continue
		}
		if s.Align {
			out = append(out, Instr{Opcode: -1, Directive: ".balign 4, 0"})
			continue
		}
		if s.Movement {
			out = append(out, Instr{Opcode: -1, Directive: ".balign 4, 0"})
		}
		ctx := &lowerCtx{labelSeq: &seq}
		// Entry-point sentinel — label only, no macro body.
		if s.Name != "" {
			ctx.append(Instr{Label: s.Name, Opcode: -1})
		}
		ctx.lowerStmts(s.Stmts)
		out = append(out, ctx.out...)
	}
	return out
}
