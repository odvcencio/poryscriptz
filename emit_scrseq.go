package poryscriptz

import (
	"strings"

	"m31labs.dev/poryscriptz/target"
)

// Emit lowers prog to a flat Instr list and renders the complete scr_seq
// assembly text for the given target.
//
// Output structure (EXACT):
//
//	<tgt.Preamble("")>
//	<tgt.Header(entryLabels)>
//	For each script:
//	  <Name>:\n
//	  For each instruction:
//	    \t<macro>\n              (zero args)
//	    \t<macro> <a1>, <a2>\n  (one or more args, joined with ", ")
func Emit(prog *Program, tgt target.GameTarget) string {
	instrs := Lower(prog)

	// Collect entry labels in order (one per script).
	var labels []string
	for _, s := range prog.Scripts {
		labels = append(labels, s.Name)
	}

	var b strings.Builder

	// File preamble (includes + .rodata).
	b.WriteString(tgt.Preamble(prog.Includes...))

	// scrdef table.
	b.WriteString(tgt.Header(labels))

	// Script bodies.
	for _, ins := range instrs {
		if ins.Directive != "" {
			b.WriteByte('\t')
			b.WriteString(ins.Directive)
			b.WriteByte('\n')
			continue
		}
		if ins.Label != "" {
			// Entry-point sentinel: emit the label line.
			b.WriteString(ins.Label)
			b.WriteString(":\n")
			// sentinel has Opcode -1 → no body line.
			if ins.Opcode == -1 {
				continue
			}
		}
		// Instruction body.
		// Macro takes precedence: control-flow wrapper macros (e.g. goto_if_unset)
		// are not base opcodes in scrcmd.json and are emitted as literal names.
		macro := ins.Macro
		if macro == "" {
			macro = tgt.MacroName(ins.Opcode)
		}
		b.WriteByte('\t')
		b.WriteString(macro)
		if len(ins.Args) > 0 {
			b.WriteByte(' ')
			b.WriteString(strings.Join(ins.Args, ", "))
		}
		b.WriteByte('\n')
	}

	return b.String()
}
