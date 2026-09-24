# poryScriptZ

Write HeartGold event scripts in `.poryz` files. Use calls, labels, and simple flow instead of writing assembly by hand. `poryz` compiles your file to the `scr_seq` macros used by the public pokeheartgold decomp. It also reads an existing `.s` file and makes an editable `.poryz` file.

This is an unofficial fan tool. It is not affiliated with Nintendo, Game Freak, or The Pokémon Company. It does not include a ROM, game art, or game text.

## Five-minute quick start

1. [Install `poryz`](#install) and put it on your `PATH`.
2. Pick one `scr_seq` `.s` file in your own decomp checkout. Make a copy or use Git so you can review your edit.
3. Turn it into `.poryz`:

   ```sh
   poryz decompile -o my-map.poryz path/to/scr_seq_0027.s
   ```

4. Open `my-map.poryz`. A command such as `End` is now `End()`. Add a message, a branch, or one of the [cookbook patterns](docs/cookbook.md). Keep the message IDs and labels that your map uses.
5. Check and compile the file:

   ```sh
   poryz fmt -w my-map.poryz
   poryz check my-map.poryz
   poryz compile -o path/to/scr_seq_0027.s my-map.poryz
   ```

Build your decomp with its normal build command. The generated `.s` file uses the decomp's own macros. The [first-script tutorial](docs/first-script.md) shows each part of a file.

## A small script

```poryz
package m

script FirstTalk {
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

`script` adds one entry to the script table. `NPCMsg(0)` uses message ID 0. Use an ID from the message file for your map. `align4()` keeps the end of the file aligned for the assembler. The file above compiles without extra setup; see [`examples/npc-dialogue.poryz`](examples/npc-dialogue.poryz) for a checked copy.

## Install

Download the `v0.3.0` archive for your system from [Releases](https://github.com/odvcencio/poryscriptz/releases/tag/v0.3.0). The archive has one `poryz` binary. Compare its SHA-256 value with `checksums.txt` on the release page.

| System | Archive | Put the binary here |
| --- | --- | --- |
| Windows x64 | `poryz_0.3.0_windows_amd64.zip` | A folder on your `PATH`; run `poryz.exe` in PowerShell |
| Windows ARM64 | `poryz_0.3.0_windows_arm64.zip` | A folder on your `PATH`; run `poryz.exe` in PowerShell |
| macOS Apple silicon | `poryz_0.3.0_darwin_arm64.tar.gz` | A folder on your `PATH`, such as `~/bin` |
| macOS Intel | `poryz_0.3.0_darwin_amd64.tar.gz` | A folder on your `PATH`, such as `~/bin` |
| Linux x64 | `poryz_0.3.0_linux_amd64.tar.gz` | A folder on your `PATH`, such as `~/bin` |
| Linux ARM64 | `poryz_0.3.0_linux_arm64.tar.gz` | A folder on your `PATH`, such as `~/bin` |

Run `poryz version` to check the install. If you use Go 1.25.1 or later, install from source:

```sh
go install github.com/odvcencio/poryscriptz/cmd/poryscriptz@v0.3.0
```

This command names the binary `poryscriptz`. Rename it to `poryz` and put it on your `PATH`. You can also build a local checkout with `GOWORK=off go build -o poryz ./cmd/poryscriptz`.

The release also has a VS Code extension package. Install it with `code --install-extension poryz-0.3.0.vsix`. Set `poryz.server.path` in VS Code if the binary is outside your `PATH`.

## Commands

| Command | What it does |
| --- | --- |
| `poryz decompile [-o file.poryz] file.s` | Turn macro assembly into editable source. |
| `poryz compile [-o file.s] file.poryz` | Write macro assembly. Without `-o`, print it. |
| `poryz check file.poryz...` | Check syntax, calls, and known label targets. |
| `poryz fmt [-w] file.poryz...` | Format source. Use `-w` to write files. |
| `poryz commands [PREFIX]` | Find script and movement commands. |
| `poryz constants [PREFIX]` | Find flag, variable, item, species, move, and map names. |
| `poryz lsp` | Start the language server for an editor. |

The binary includes its vocabulary. If you work with a fork that changes `scrcmd.json`, pass `--scrcmd path/to/scrcmd.json`. Put that fork's `macros.json` and `decomp_macros.json` in the same folder.

The [helpful patterns](docs/patterns.md) cover messages, gifts, trainer battles, flags, movement, and warps. The [command reference](docs/command-reference.md) is generated from macro declarations in the public decomp. The [cookbook](docs/cookbook.md) has 26 complete examples. The [migration guide](docs/migration.md) helps you move a map from `.s` to `.poryz`.

## Editor support

The VS Code extension starts `poryz lsp`. The Go server checks a file as you type, shows command signatures on hover, jumps to labels, completes constants, and sends syntax colors from the embedded gotreesitter grammar. The extension is a client; parsing and analysis run in Go. Other editors can use `poryz lsp` over standard input and output. The grammar query is in [`queries/highlights.scm`](queries/highlights.scm).

## Compatibility and proof

The decompiler keeps the command order, table entries, labels, alignment, and includes needed for a binary round trip. A local parity tool assembles both the original `.s` and the recompiled `.s`, then compares their bytes. It passed for all 965 public `scr_seq` files and all 35 scripts in a separate edited test corpus at the time of v0.3.0. See [round-trip testing](docs/round-trip.md) for the commands and limits.

The project is [MIT licensed](LICENSE).
