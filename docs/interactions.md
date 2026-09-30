# NPC interactions

Build this branch from source to use these additive HGSS helpers. The published
v0.3.0 binary predates them. Existing source and raw scripts keep their behavior.

## A readable greeting

This raw NPC talk sequence:

```poryz
package m
script Greeting {
    LockAll()
    FacePlayer()
    NPCMsg(0)
    WaitButton()
    CloseMsg()
    ReleaseAll()
    End()
}
align4()
```

can be written as:

```poryz
package m
interaction Greeting {
    dialogue(0)
}
align4()
```

The entry table, message ID, command order, and assembled bytes are the same.
`interaction` is for an NPC the player talks to. It emits LockAll and FacePlayer
before its body, then ReleaseAll and End. `dialogue` emits NPCMsg, WaitButton,
and CloseMsg. It also works inside a raw script, where you own the locks.
`say(id)` keeps its earlier NPCMsg + CloseMsg behavior, with no wait.

Use message IDs already present in the map's bank. `dialogue(MSG_GREETING)` and
`if ask(msg_offer)` accept existing symbolic references from an included `.h`
file. They emit those symbols for the assembler to resolve. The compiler checks
literal IDs fit NPCMsg's 8-bit field and checks symbol spelling, but does not read
header definitions or prove their numeric values. These helpers do not write
message banks or accept inline English text. No alias declarations are added.

## Ask without numeric menu results

```poryz
package m
interaction Offer {
    if flag(FLAG_UNK_042) {
        dialogue(3)
    } else if ask(0) {
        dialogue(1)
        flag_set(FLAG_UNK_042)
    } else {
        dialogue(2)
    }
}
align4()
```

Choose a flag owned by your event and replace the four message IDs with your
bank's repeat, offer, accepted, and declined messages. The first Yes completes
this conversation; a repeat visit takes the first branch. No and B cancellation
take the final branch and leave the completion flag clear. The example grants
no item; item rewards still use the existing checked three-argument `give_item`
helper in a raw script with an explicit bag-full branch.

`ask(message)` emits NPCMsg, YesNo(VAR_SPECIAL_RESULT), CloseMsg, and a comparison
with **0**. HGSS returns 0 for Yes and 1 for No, including cancellation. `ask`
is supported only as an `if` predicate, and clobbers `VAR_SPECIAL_RESULT`.
Nested asks are safe because each branch tests the result before its body runs.
Do not expect an earlier scratch result to survive an ask.

The polarity follows the public engine's
[sub_020416E4](https://github.com/pret/pokeheartgold/blob/9d8b7591f09b65804da2fb2dfd56f320633e0d36/src/scrcmd_c.c)
and its Yes/No menu handler. Raw YesNo remains available for manual result use.

## Cleanup and the escape hatch

An interaction body initially accepts `dialogue`, `flag_set`, `flag_clear`,
`var_set`, `var_copy`, and structured `if`/`else` conditions using `flag`, `var`,
or `ask`. Each accepted path falls through to the generated release and end.
Raw commands, early exits, manual locks, movement, warps, gifts with external
branches, loops, and switches are rejected with their source line and column.
Use a raw `script` when your event needs them and manage its lifecycle explicitly.

Unsupported inherited Go forms such as assignment, return, break, functions,
and imports now produce diagnostics rather than disappearing from compiled
output. `AsmGoto` and `AsmReturn` are still explicit raw-script escape hatches.

## Migrate one event and reuse it

Keep the decompiler's direct calls as your baseline. Replace only a complete
NPC talk entry with `interaction`, keeping its name, table position, includes,
message IDs, flags, and alignment. Replace display/wait/close runs with `dialogue`.
For Yes/No, verify the original branch really means Yes before converting it to
`if ask`. Build with the same macros and compare assembled bytes for unchanged
behavior; then play Yes, No, cancellation, and repeat visits.

These helpers work across maps without new opcodes or per-map compiler changes.
They do not automatically rewrite imported scripts or update a game's compiler
pin. See [round-trip testing](round-trip.md) for the corpus gate and the opt-in
managed-flow binary tests.
