package target

import "m31labs.dev/poryscriptz/vocab"

type GameTarget interface {
	Vocabulary() *vocab.Table
	MacroName(opcode int) string // script.inc macro name, or "scrcmd_NNN"
	// Preamble emits the file header every real scr_seq/*.s begins with: the
	// #include chain, `.include "asm/macros/script.inc"`, and `.rodata`.
	// eventHeader is the per-file include or "" for synthetic tests.
	Preamble(eventHeader string) string
	Header(labels []string) string // scrdef table: one "scrdef <label>" per entry + scrdef_end
}
