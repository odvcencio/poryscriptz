# Move from `.s` to `.poryz`

Start with one map. Keep your original `.s` file under Git so you can compare the output.

1. Run `poryz decompile -o map.poryz path/to/map.s`.
2. Open `map.poryz`. Each macro is a call. For example, `NPCMsg 0` becomes `NPCMsg(0)`. A script table entry appears in `entries(...)`; its body is a `script` block. Other labels are `label` blocks.
3. Run `poryz check map.poryz` and `poryz fmt -w map.poryz`.
4. Run `poryz compile -o path/to/map.s map.poryz`.
5. Review the `.s` diff and build the decomp with its normal build command.

The decompiler keeps direct command calls so the first edit is easy to review. You can then replace common command runs with helpers such as `give_item`, `move_actor_and_wait`, or `warp_to`. Keep the same message IDs and event constants.

A file with `dialect pret` uses the public decomp's macro names. A file with `dialect legacy` uses the older lower case names. The decompiler writes this line for you. `header init` holds map init commands. `table_end()` keeps an extra script table end marker if the original file has one. `source "path.s"` keeps a source include in place. These forms preserve bytes during a round trip.

The decompiler stops with a file position and a fix hint if it finds an unknown macro or an assembly directive it cannot express. Check that you used the macro definitions from the same decomp revision. You can regenerate the vocabulary with `GOWORK=off go run ./cmd/genvocab /path/to/decomp` when you build poryz from source.
