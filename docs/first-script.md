# Your first script

This tutorial starts with a small NPC talk script. You can read and change it without knowing assembly.

## 1. Make a source file

Save this as `hello.poryz`:

```poryz
package m

script Hello {
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

`package m` marks the file as poryz source. `script Hello` adds a script table entry named `Hello`. Each call runs one command. The braces show where the script starts and ends.

`LockAll()` keeps the player in place. `FacePlayer()` turns the NPC. `NPCMsg(0)` opens message 0 from the map's message file. `WaitButton()` waits for a button. `CloseMsg()` closes the box. `ReleaseAll()` gives control back. `End()` stops the script.

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

script Hello {
    if flag(FLAG_UNK_042) {
        NPCMsg(1)
    } else {
        NPCMsg(0)
        SetFlag(FLAG_UNK_042)
    }
    WaitButton()
    CloseMsg()
    End()
}
```

Replace `FLAG_UNK_042` with a flag that belongs to your event. Use `poryz constants FLAG_` to find names, and check the map's event header for local flags. See the [flag gate recipe](cookbook.md#open-an-event-with-a-flag) for a complete file.

If `poryz check` reports `file:line:col`, go to that call. The error tells you what the compiler expected. The [command reference](command-reference.md) shows argument names for each macro.
