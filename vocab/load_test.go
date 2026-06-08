package vocab

import "testing"

func TestLoadScrcmdJSON(t *testing.T) {
	tbl, err := LoadScrcmdJSON("testdata/scrcmd.json")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	cmd, ok := tbl.ByName("end")
	if !ok {
		t.Fatal(`command "end" not found`)
	}
	if len(cmd.Args) != 0 {
		t.Fatalf("end: want 0 args, got %d", len(cmd.Args))
	}
	if _, ok := tbl.ByName("setflag"); !ok {
		t.Fatal(`command "setflag" not found`)
	}
}
