//go:build decompparity

package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	poryz "github.com/odvcencio/poryscriptz"
	"github.com/odvcencio/poryscriptz/target"
)

// Original synthetic flows contain no extracted dialogue or game assets.
// Both sides use the same local assembler, macros, and constants.
func TestManagedFlowBytes(t *testing.T) {
	root, mwas := os.Getenv("POKEHG_ROOT"), os.Getenv("POKEHG_MWAS")
	if root == "" || mwas == "" {
		t.Skip("set POKEHG_ROOT to an isolated public decomp and POKEHG_MWAS to an authorized local assembler")
	}
	for _, tool := range []string{"wine", "arm-none-eabi-objcopy"} {
		if _, err := exec.LookPath(tool); err != nil {
			t.Fatalf("configured parity test requires %s: %v", tool, err)
		}
	}
	tgt, err := target.HGSS("")
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ name, managed, raw string }{
		{"greeting",
			"interaction Greeting { dialogue(7) }",
			"script Greeting { LockAll(); FacePlayer(); NPCMsg(7); WaitButton(); CloseMsg(); ReleaseAll(); End() }"},
		{"yes-no-offer",
			"interaction Offer { if ask(7) { dialogue(8) } else { dialogue(9) } }",
			`script Offer {
    LockAll(); FacePlayer(); NPCMsg(7); YesNo(VAR_SPECIAL_RESULT); CloseMsg()
    Compare(VAR_SPECIAL_RESULT, 0); GoToIfNe(_Declined)
    NPCMsg(8); WaitButton(); CloseMsg(); GoTo(_Finished)
}
label _Declined { NPCMsg(9); WaitButton(); CloseMsg() }
label _Finished { ReleaseAll(); End() }`},
		{"no-else",
			"interaction Offer { if ask(7) { dialogue(8) } }",
			`script Offer {
    LockAll(); FacePlayer(); NPCMsg(7); YesNo(VAR_SPECIAL_RESULT); CloseMsg()
    Compare(VAR_SPECIAL_RESULT, 0); GoToIfNe(_Finished)
    NPCMsg(8); WaitButton(); CloseMsg()
}
label _Finished { ReleaseAll(); End() }`},
		{"symbolic-message",
			"interaction Greeting { dialogue(MSG_GREETING) }",
			"script Greeting { LockAll(); FacePlayer(); NPCMsg(MSG_GREETING); WaitButton(); CloseMsg(); ReleaseAll(); End() }"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tmp := t.TempDir()
			var binaries [][]byte
			for _, side := range []struct{ name, source string }{{"managed", tc.managed}, {"raw", tc.raw}} {
				asm, ds, err := poryz.Compile([]byte("package m\n"+side.source+"\nalign4()\n"), tgt)
				if err != nil || len(ds) != 0 {
					t.Fatalf("%s compile: %v %v", side.name, err, ds)
				}
				// An original symbolic message reference; no message bank or text.
				asm = "#define MSG_GREETING 7\n" + asm
				file := filepath.Join(tmp, side.name+".s")
				if err := os.WriteFile(file, []byte(asm), 0644); err != nil {
					t.Fatal(err)
				}
				binary, err := assemble(file, tmp, side.name, root, mwas, "")
				if err != nil {
					t.Fatalf("%s assembly: %v", side.name, err)
				}
				binaries = append(binaries, binary)
			}
			if !bytes.Equal(binaries[0], binaries[1]) {
				t.Fatalf("binary mismatch: managed=%x raw=%x", binaries[0], binaries[1])
			}
			t.Logf("%s: %d byte-identical bytes", tc.name, len(binaries[0]))
		})
	}
}
