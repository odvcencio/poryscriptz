# poryScriptZ

poryScriptZ compiles `.poryz` event scripts to HeartGold `scr_seq` assembly.
It emits the existing `asm/macros/script.inc` commands. It does not add VM
opcodes.

Run the compiler from this repository:

```sh
go run ./cmd/poryscriptz --scrcmd vocab/testdata/scrcmd.json -o event.s event.poryz
```

Forge passes its HeartGold vocabulary with `--scrcmd`. The compiler loads
`macros.json` from the same directory as `scrcmd.json`.

## Source layout

```go
package m
include "fielddata/script/scr_seq/event_D02R0103.h"

script FossilGift {
    lockall()
    faceplayer()
    give_item(ITEM_DOME_FOSSIL, 1, _BagFull)
    flag_set(FLAG_HIDE_ITEMBALL_D02R0103_FOSSILS)
    closemsg()
    releaseall()
    end()
}

label _BagFull {
    callstd(std_bag_is_full)
    closemsg()
    releaseall()
    end()
}
align4()
```

Each `script` adds one entry to the `scrdef` table. Each `label` marks a
branch target without an entry. Empty, adjacent labels share an address.
`align4()` emits `.balign 4, 0` at its source position. Use `include` for
relative `.h` files that define event or message symbols. The compiler
rejects absolute paths and paths with `..` traversal.

Use a `movement` block for a movement path:

```go
script WalkOut {
    move_actor_and_wait(0, _Walk)
    end()
}

movement _Walk {
    WalkNormalNorth(2)
    EndMovement()
}
```

The compiler aligns each movement block to four bytes. It checks the listed
movement macros and requires `EndMovement()` as the last command. The v0.3
movement list covers normal, fast, slightly fast, and on the spot walking in
four directions. It also covers far jumps, fast on the spot jumps, `Delay16`,
`EmoteExclamationMark`, and `SetVisible`.

## Event patterns

These calls expand to the listed assembly commands. Arguments pass through to
the assembler after the compiler checks their shape.

| Source call | Emitted commands |
| --- | --- |
| `message(id)` | `npc_msg id` |
| `say(id)` | `npc_msg id`; `closemsg` |
| `trainer_battle_simple(trainer)` | `trainer_battle trainer, 0, 0, 0` |
| `give_item(item, count, bagFullLabel)` | `goto_if_no_item_space item, count, bagFullLabel`; `callstd std_give_item_verbose` |
| `move_actor(actor, path)` | `apply_movement actor, path` |
| `move_actor_and_wait(actor, path)` | `apply_movement actor, path`; `wait_movement` |
| `flag_set(flag)` / `flag_clear(flag)` | `setflag flag` / `clearflag flag` |
| `var_set(variable, value)` | `setvar variable, value` |
| `var_copy(destination, source)` | `copyvar destination, source` |
| `warp_to(map, warp, x, y, direction)` | `warp map, warp, x, y, direction` |
| `if_flag_set(flag, label)` / `if_flag_unset(flag, label)` | `goto_if_set flag, label` / `goto_if_unset flag, label` |
| `if_var_eq(variable, value, label)` | `compare_var_to_value variable, value`; `goto_if_eq label` |
| `if_var_ne(variable, value, label)` | `compare_var_to_value variable, value`; `goto_if_ne label` |

`give_item` checks Bag space and calls the standard gift routine. Add the
hide flag, message close, and release commands your event needs. A trainer
battle does not check the result. Add a result check before you grant a win.
`move_actor` starts movement without waiting. Use `move_actor_and_wait` when
the next command must run after movement ends. Point the path at a `movement`
block in the same file. User labels must not use `_L` followed only by digits.
The compiler reserves those names for generated flow labels.
`var_set` and `if_var_eq` or `if_var_ne` take literal values or constants.
Use `var_copy` to copy one variable into another. Branch patterns require a
defined `script` or `label` target. Movement patterns require a defined
`movement` target.

`give_item` accepts an item or quantity from a `VAR_*` source. Numeric item
and quantity literals must be at most `0x3FFF`. The HGSS `item_vars` macro
reads values of `0x4000` or more as variable IDs. Numeric message IDs must
fit one byte. Other numeric pattern operands and movement repeat counts must
fit a 16-bit field. The compiler leaves symbolic constants for the assembler
to resolve.

You can also call the existing HGSS command names, use `scrcmd_NNN(...)` for
an unnamed opcode, and use structured `if`, `for`, and `switch` flow. The
v0.3 patterns keep the command order and argument text visible in output.

## Checks

```sh
go test ./...
```

The golden tests pin command lowering. A byte parity gate must assemble the
old `.s` and new `.poryz` output with the same decomp toolchain, then compare
the assembled member bytes. Text equality alone does not prove byte parity.
Set `POKEHG_ROOT` to an isolated decomp checkout when you run tagged parity
tests. Set `POKEHG_MWAS` if the assembler is outside that checkout. Tagged
tests skip when `POKEHG_ROOT` is absent.
