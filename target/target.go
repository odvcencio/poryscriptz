package target

import "m31labs.dev/poryscriptz/vocab"

type GameTarget interface {
	Name() string // human-readable target name, e.g. "HGSS"
	Vocabulary() *vocab.Table
	// Macros returns the script.inc macro table (names + arg shapes).
	// Macro calls in .poryz are emitted verbatim; the decomp assembler expands them.
	Macros() *vocab.MacroTable
	MacroName(opcode int) string // script.inc macro name, or "scrcmd_NNN"
	// Preamble emits the file header every real scr_seq/*.s begins with: the
	// #include chain, `.include "asm/macros/script.inc"`, and `.rodata`.
	// eventHeader is the per-file include or "" for synthetic tests.
	Preamble(eventHeader string) string
	Header(labels []string) string // scrdef table: one "scrdef <label>" per entry + scrdef_end
}
