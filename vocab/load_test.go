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

func TestLoadMacrosJSON(t *testing.T) {
	mt, err := LoadMacrosJSON("testdata/macros.json")
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	// Zero-arg macros.
	for _, name := range []string{"closemsg", "faceplayer", "lockall", "releaseall", "wait_movement"} {
		e, ok := mt.ByName(name)
		if !ok {
			t.Fatalf("macro %q not found", name)
		}
		if len(e.Args) != 0 {
			t.Fatalf("%s: want 0 args, got %d", name, len(e.Args))
		}
	}

	// trainer_battle: 4 args.
	e, ok := mt.ByName("trainer_battle")
	if !ok {
		t.Fatal(`macro "trainer_battle" not found`)
	}
	if len(e.Args) != 4 {
		t.Fatalf("trainer_battle: want 4 args, got %d", len(e.Args))
	}

	// give_mon: 6 args.
	e, ok = mt.ByName("give_mon")
	if !ok {
		t.Fatal(`macro "give_mon" not found`)
	}
	if len(e.Args) != 6 {
		t.Fatalf("give_mon: want 6 args, got %d", len(e.Args))
	}

	// warp: 5 args.
	e, ok = mt.ByName("warp")
	if !ok {
		t.Fatal(`macro "warp" not found`)
	}
	if len(e.Args) != 5 {
		t.Fatalf("warp: want 5 args, got %d", len(e.Args))
	}

	// apply_movement: 2 args.
	e, ok = mt.ByName("apply_movement")
	if !ok {
		t.Fatal(`macro "apply_movement" not found`)
	}
	if len(e.Args) != 2 {
		t.Fatalf("apply_movement: want 2 args, got %d", len(e.Args))
	}

	// check_badge: 2 args.
	_, ok = mt.ByName("check_badge")
	if !ok {
		t.Fatal(`macro "check_badge" not found`)
	}

	// Unknown macro must not be found.
	if _, ok := mt.ByName("definitely_not_a_macro"); ok {
		t.Fatal("unexpected macro found for unknown name")
	}
}
