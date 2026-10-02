# Your first script

This tutorial starts with a small NPC talk script. You can read and change it without knowing assembly.

## 1. Make a source file

Save this as `hello.poryz`:

```poryz
package m

interaction Hello {
    dialogue(0)
}
align4()
```

Build this branch from source for the interaction helpers; the published v0.3.0 binary does not have them yet.

`package m` marks the file as poryz source. `interaction Hello` adds an NPC talk entry named `Hello`, keeps actors in place, and faces the player. `dialogue(0)` opens message 0 from the map's message bank, waits for a button, and closes the box. The interaction then gives control back and stops. Use a raw `script` for other kinds of map triggers or manual command sequences.

## 2. Check and compile it

```sh
poryz check hello.poryz
poryz compile -o hello.s hello.poryz
```

Open `hello.s` to see the macro assembly. The output keeps one command per line, so you can review it. The `ScrDef` line points to `Hello`. The final `.balign` line comes from `align4()`.

## 3. Use it in a map

A real map has its own message IDs, event header, and script table. The safe way to start is to decompile the map's current `.s` file:

```sh
poryz decompile -o my-map.poryz path/to/my-map.s
```

Find a `script` block and change its calls. Keep the map's `include` lines. Use the real message ID in place of `0`. Run `poryz check`, then compile back to the map's `.s` file. Build the decomp as usual.

## 4. Add a branch

Use a flag when the NPC should say something else after an event:

```poryz
package m

interaction Hello {
    if flag(FLAG_UNK_042) {
        dialogue(1)
    } else {
        dialogue(0)
        flag_set(FLAG_UNK_042)
    }
}
```

Replace `FLAG_UNK_042` with a flag that belongs to your event. Use `poryz constants FLAG_` to find names, and check the map's event header for local flags. See the [flag gate recipe](cookbook.md#open-an-event-with-a-flag) for a complete file.

If `poryz check` reports `file:line:col`, go to that call. The error tells you what the compiler expected. The [command reference](command-reference.md) shows argument names for each macro.
