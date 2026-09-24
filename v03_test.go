package poryscriptz

import (
	"os"
	"strings"
	"testing"
)

func TestV03PatternsGolden(t *testing.T) {
	src, err := os.ReadFile("testdata/v03_patterns.poryz")
	if err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile("testdata/v03_patterns.golden")
	if err != nil {
		t.Fatal(err)
	}
	got, diags, err := Compile(src, hgssT(t))
	if err != nil || len(diags) != 0 {
		t.Fatalf("compile: err=%v diags=%v", err, diags)
	}
	if got != string(want) {
		t.Fatalf("golden mismatch:\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
	if strings.Count(got, "\tscrdef ") != 1 {
		t.Fatalf("non-entry labels reached scrdef table:\n%s", got)
	}
}

func TestV03PatternDiagnostics(t *testing.T) {
	cases := []struct {
		name, call, want string
	}{
		{"arity", "give_item(ITEM_POTION, 1)", "expects 3 args, got 2"},
		{"flag type", "flag_set(ITEM_POTION)", "requires FLAG_ identifier"},
		{"var type", "var_set(FLAG_X, 1)", "requires VAR_ identifier"},
		{"label type", "give_item(ITEM_POTION, 1, 3)", "requires label identifier"},
		{"map type", "warp_to(ITEM_POTION, 0, 8, 11, DIR_SOUTH)", "requires map id"},
		{"message type", "message(FLAG_X)", "requires message id"},
		{"compare var value", "if_var_eq(VAR_X, VAR_Y, _Done)", "not VAR_"},
		{"set var value", "var_set(VAR_X, VAR_Y)", "not VAR_"},
		{"gift quantity", "give_item(ITEM_POTION, FLAG_COUNT, _Full)", "requires gift quantity"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			src := []byte("package m\nscript S {\n" + tc.call + "\n}\n")
			_, diags, err := Compile(src, hgssT(t))
			if err != nil {
				t.Fatal(err)
			}
			if len(diags) != 1 || !strings.Contains(diags[0].Msg, tc.want) {
				t.Fatalf("want %q, got %v", tc.want, diags)
			}
			if diags[0].Line < 1 || diags[0].Col < 1 {
				t.Fatalf("missing source position: %v", diags[0])
			}
		})
	}
}

func TestV03UnsafeIncludeAndDuplicateLabel(t *testing.T) {
	src := []byte("package m\ninclude \"../outside.h\"\nscript S {}\nlabel S {}\n")
	_, diags, err := Compile(src, hgssT(t))
	if err != nil {
		t.Fatal(err)
	}
	if len(diags) != 2 || !strings.Contains(diags[0].Msg, "relative .h path") || !strings.Contains(diags[1].Msg, "duplicate label") {
		t.Fatalf("want include and label diagnostics, got %v", diags)
	}
}

func TestV03MovementDiagnostics(t *testing.T) {
	cases := []struct {
		name, body, want string
	}{
		{"missing end", "WalkNormalNorth(2)", "requires EndMovement"},
		{"unknown macro", "WalkUnknownNorth()\nEndMovement()", "unknown movement macro"},
		{"wrong count", "WalkNormalNorth(FLAG_X)\nEndMovement()", "repeat count must be an integer"},
		{"command after end", "EndMovement()\nWalkNormalNorth()", "must be the last"},
		{"repeat overflow", "WalkNormalNorth(65536)\nEndMovement()", "16-bit field"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			src := []byte("package m\nscript S { end() }\nmovement Path {\n" + tc.body + "\n}\n")
			_, diags, err := Compile(src, hgssT(t))
			if err != nil {
				t.Fatal(err)
			}
			if len(diags) == 0 || !strings.Contains(diags[0].Msg, tc.want) {
				t.Fatalf("want %q, got %v", tc.want, diags)
			}
		})
	}
}

func TestV03PatternTargetDiagnostics(t *testing.T) {
	cases := []struct {
		name, source, want string
	}{
		{"missing gift branch", "script S { give_item(ITEM_POTION, 1, _Missing) }", "defined branch label"},
		{"missing flag branch", "script S { if_flag_set(FLAG_GAME_CLEAR, _Missing) }", "defined branch label"},
		{"nested var branch", "script S { if flag(FLAG_GAME_CLEAR) { if_var_eq(VAR_X, 1, _Missing) } }", "defined branch label"},
		{"branch into movement", "script S { if_flag_unset(FLAG_GAME_CLEAR, _Move) }\nmovement _Move { EndMovement() }", "defined branch label"},
		{"missing movement", "script S { move_actor(0, _Missing) }", "defined movement section"},
		{"non-movement target", "script S { move_actor_and_wait(0, _Label) }\nlabel _Label {}", "defined movement section"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, diags, err := Compile([]byte("package m\n"+tc.source+"\n"), hgssT(t))
			if err != nil {
				t.Fatal(err)
			}
			if len(diags) != 1 || !strings.Contains(diags[0].Msg, tc.want) {
				t.Fatalf("want %q, got %v", tc.want, diags)
			}
			if diags[0].Line < 1 || diags[0].Col < 1 {
				t.Fatalf("missing target position: %v", diags[0])
			}
		})
	}
}

func TestV03PatternNumericBounds(t *testing.T) {
	cases := []struct {
		name, call, want string
	}{
		{"message byte", "message(256)", "<= 255"},
		{"say byte", "say(0x100)", "<= 255"},
		{"trainer short", "trainer_battle_simple(65536)", "<= 65535"},
		{"gift item literal", "give_item(0x4000, 1, _Done)", "<= 16383"},
		{"gift quantity literal", "give_item(ITEM_POTION, 0x4000, _Done)", "<= 16383"},
		{"actor short", "move_actor(65536, _Move)", "<= 65535"},
		{"var value short", "var_set(VAR_X, 65536)", "<= 65535"},
		{"compare short", "if_var_ne(VAR_X, 65536, _Done)", "<= 65535"},
		{"warp short", "warp_to(MAP_PALLET, 0, 65536, 0, 0)", "<= 65535"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			src := []byte("package m\nscript S { " + tc.call + " }\nlabel _Done {}\nmovement _Move { EndMovement() }\n")
			_, diags, err := Compile(src, hgssT(t))
			if err != nil {
				t.Fatal(err)
			}
			if len(diags) != 1 || !strings.Contains(diags[0].Msg, tc.want) {
				t.Fatalf("want %q, got %v", tc.want, diags)
			}
		})
	}
}

func TestV03PatternNumericEdges(t *testing.T) {
	src := []byte("package m\nscript S { message(255); trainer_battle_simple(65535); give_item(0x3FFF, 0x3FFF, _Done); move_actor(65535, _Move); var_set(VAR_X, 65535); if_var_eq(VAR_X, 65535, _Done); warp_to(MAP_PALLET, 65535, 65535, 65535, 65535) }\nlabel _Done {}\nmovement _Move { WalkNormalNorth(65535); EndMovement() }\n")
	_, diags, err := Compile(src, hgssT(t))
	if err != nil || len(diags) > 0 {
		t.Fatalf("valid width edges must compile: err=%v diags=%v", err, diags)
	}
}

func TestV03GiftVariableQuantity(t *testing.T) {
	src := []byte("package m\nscript S { give_item(VAR_SPECIAL_x8004, VAR_SPECIAL_x8005, _Full) }\nlabel _Full { end() }\n")
	asm, diags, err := Compile(src, hgssT(t))
	if err != nil || len(diags) > 0 {
		t.Fatalf("dynamic gift must compile: err=%v diags=%v", err, diags)
	}
	if !strings.Contains(asm, "goto_if_no_item_space VAR_SPECIAL_x8004, VAR_SPECIAL_x8005, _Full") {
		t.Fatalf("dynamic gift changed macro arguments:\n%s", asm)
	}
}

func TestV03GeneratedLabelNamespaceReserved(t *testing.T) {
	src := []byte("package m\nscript S { if flag(FLAG_X) { end() } }\nlabel _L0 {}\n")
	_, diags, err := Compile(src, hgssT(t))
	if err != nil {
		t.Fatal(err)
	}
	if len(diags) != 1 || !strings.Contains(diags[0].Msg, "reserved generated-label namespace") {
		t.Fatalf("want reserved-label diagnostic, got %v", diags)
	}
}
