package poryscriptz

import (
	"strings"

	"github.com/odvcencio/poryscriptz/target"
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
	if prog.Entries != nil {
		labels = prog.Entries
	}

	var b strings.Builder

	// File preamble (includes + .rodata).
	b.WriteString(tgt.Preamble(prog.Includes...))

	// scrdef table.
	if !prog.HasInitHeader {
		if prog.PretDialect {
			for _, l := range labels {
				b.WriteString("\tScrDef ")
				b.WriteString(l)
				b.WriteByte('\n')
			}
			b.WriteString("\tScrDefEnd\n")
		} else {
			b.WriteString(tgt.Header(labels))
		}
		for i := 0; i < prog.ExtraTableEnds; i++ {
			if prog.PretDialect {
				b.WriteString("\tScrDefEnd\n")
			} else {
				b.WriteString("\tscrdef_end\n")
			}
		}
	} else {
		for _, ins := range Lower(&Program{Sections: []Section{{Stmts: prog.InitHeader}}}) {
			name := ins.Macro
			if name == "" {
				if prog.PretDialect {
					if canonical, ok := tgt.(interface{ CanonicalOpcode(int) string }); ok {
						name = canonical.CanonicalOpcode(ins.Opcode)
					}
				}
				if name == "" {
					name = tgt.MacroName(ins.Opcode)
				}
			}
			if prog.PretDialect {
				if canonical, ok := tgt.(interface{ CanonicalMacro(string) string }); ok {
					name = canonical.CanonicalMacro(name)
				}
			}
			b.WriteByte('\t')
			b.WriteString(name)
			if len(ins.Args) > 0 {
				b.WriteByte(' ')
				b.WriteString(strings.Join(ins.Args, ", "))
			}
			b.WriteByte('\n')
		}
	}

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
			if prog.PretDialect {
				if canonical, ok := tgt.(interface{ CanonicalOpcode(int) string }); ok {
					macro = canonical.CanonicalOpcode(ins.Opcode)
				}
			}
			if macro == "" {
				macro = tgt.MacroName(ins.Opcode)
			}
		}
		if prog.PretDialect {
			if canonical, ok := tgt.(interface{ CanonicalMacro(string) string }); ok {
				macro = canonical.CanonicalMacro(macro)
			}
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
