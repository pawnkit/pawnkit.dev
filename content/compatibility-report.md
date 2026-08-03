# Compatibility report

The August 3 signed set passed the PawnKit workflow on Linux, Windows, and
macOS. It is the current set used by PawnKit Actions and editor tooling.

## Versions

| Tool | Version |
| --- | --- |
| PawnKit CLI | 1.34.21 |
| Pawn formatter | 1.4.10 |
| Pawn linter | 1.8.38 |
| Pawn language server | 0.34.18 |
| Pawn test runner | 1.2.11 |

The exact commits, release archives, sizes, and SHA-256 hashes are in the
[release set](/release-sets/toolchain-signed-2026-08-03-9.json).

## What passed

The same release set ran against the small SA-MP and open.mp projects in
`pawn-corpus`. `pawn doctor`, `pawn fmt --check`, and `pawn lint` passed on all
three platforms. Linux also built all three golden projects, ran the tests,
and exercised the packaged run adapter.

The compiler-backed build used Pawn 3.10.10 on Linux. The run check used the
published adapter protocol fixture; it did not start a live open.mp server.

See the [successful smoke and verification run](https://github.com/pawnkit/pawn-actions/actions/runs/30786284529)
for the recorded jobs and supply-chain checks.

## Known limits

- The language server was installed and its version checked. This run did not
  exercise the LSP protocol.
- Native plugins and debugging remain experimental.
- Real-server integration is not part of this signed set.
- The small projects do not replace the larger corpus and real-project suites
  used by individual tools.

Use the versions above together. Mixing older tools may produce different
project discovery, diagnostics, or protocol behavior.
