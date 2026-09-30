# Round-trip test

The `roundtrip` tool reads each source `.s`, decompiles it, compiles the `.poryz`, assembles both `.s` files with the same MWAS and headers, and compares the binary bytes after `objcopy`.

| Corpus | Files | Byte-identical | Failures |
| --- | ---: | ---: | ---: |
| Public pret `scr_seq` checkout | 965 | 965 | 0 |
| Edited external script corpus | 35 | 35 | 0 |
| Total | 1,000 | 1,000 | 0 |

The public checkout used commit `9d8b7591f09b65804da2fb2dfd56f320633e0d36`. The test used the local MWAS 2.0/sp2p2 assembler and `arm-none-eabi-objcopy`. Message ID headers were generated from the public `.gmm` row IDs in the scratch checkout. The edited corpus used its matching decomp headers. No message text or binary output is stored here.

Run the test after you provide your own MWAS and generated message headers:

```sh
GOWORK=off go run ./cmd/roundtrip \
  -decomp /path/to/public-decomp \
  -mwas /path/to/mwasmarm.exe \
  -headers /path/to/generated-files
```

Add `-hack /path/to/edited/scr_seq -hack-root /path/to/matching-decomp` to test edited scripts. The test prints every failure by class and exits nonzero for a mismatch, assembly failure, or an empty public or explicitly requested edited corpus. Invalid options, including a worker count below one, exit with status 2.

The old corpus gate proves decompile/compile preservation. Managed interactions
also have parser, resolver, branch-effect, cleanup, diagnostic, formatter, and
determinism tests. Their opt-in binary tests assemble original synthetic raw
flows and the new syntax using the same authorized local toolchain:

```sh
POKEHG_ROOT=/path/to/isolated/public-decomp \
POKEHG_MWAS=/path/to/mwasmarm.exe \
WINEPREFIX=/path/to/test-prefix WINEARCH=win32 \
GOWORK=off go test -tags decompparity ./cmd/roundtrip -run TestManagedFlowBytes -v
```

No ROM build is needed. No extracted game binaries or message text are stored
in these tests. Menu branch-effect tests use the engine's verified result values;
they complement byte equivalence and do not claim an emulator playtest.
