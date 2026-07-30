# Compatibility report

The July 30 signed set passed the PawnKit workflow on Linux, Windows, and
macOS. It is the current set used by PawnKit Actions and editor tooling.

## Versions

| Tool | Version |
| --- | --- |
| PawnKit CLI | 1.6.0 |
| Pawn formatter | 1.4.4 |
| Pawn linter | 1.8.7 |
| Pawn language server | 0.33.71 |
| Pawn test runner | 1.2.5 |

The exact commits, release archives, sizes, and SHA-256 hashes are in the
[release set](/release-sets/toolchain-signed-2026-07-30.json).

## What passed

The same release set ran against the small SA-MP and open.mp projects in
`pawn-corpus`. `pawn doctor`, `pawn fmt --check`, and `pawn lint` passed on all
three platforms. Linux also ran the tests, compiler build, and packaged run
adapter.

The compiler-backed build used Pawn 3.10.10 on Linux. The run check used the
published adapter protocol fixture; it did not start a live open.mp server.

See the [successful Actions run](https://github.com/pawnkit/pawn-actions/actions/runs/30516012521)
for the recorded jobs.

## Known limits

- The language server was installed and its version checked. This run did not
  exercise the LSP protocol.
- Native plugins and debugging remain experimental.
- Real-server integration is not part of this signed set.
- The small projects do not replace the larger corpus and real-project suites
  used by individual tools.

Use the versions above together. Mixing older tools may produce different
project discovery, diagnostics, or protocol behavior.
