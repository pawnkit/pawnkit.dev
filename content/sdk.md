# PawnKit SDK

The PawnKit SDK is a set of Go modules for building Pawn tools. Pick the module
that owns the data you need instead of rebuilding that layer in your project.

## Start with the owning module

| Need | Module |
| --- | --- |
| Source ranges, edits, diagnostics, and protocol values | `pawnkit-core` |
| Tokens and lossless syntax trees | `pawn-parser` |
| Preprocessing, symbols, references, tags, and control flow | `pawn-analysis` |
| Project roots, manifests, profiles, includes, and toolchains | `pawn-project` |
| SA-MP and open.mp API facts | `pawn-api` |
| Bounded AMX loading and execution | `goamx` |

Parsing only describes source structure. Use `pawn-analysis` when macros,
active conditional branches, includes, or symbol meaning matter. Use
`pawn-project` before analysis so every tool resolves the same root, profile,
defines, and include paths.

## Add a module

Install the released version you intend to support:

```sh
go get github.com/pawnkit/pawn-parser@latest
```

For reproducible builds, replace `latest` with a tag and commit `go.mod` and
`go.sum`. Do not use a local `replace` directive in a published module.

Each repository documents its public packages and maturity. Preview modules may
change before 1.0, so keep dependency updates deliberate and test them against
your own fixtures.

## Keep the layers consistent

Do not maintain a private copy of project discovery, API data, preprocessor
rules, or diagnostics. Fix shared facts in their owning module, then update the
consumer to the released tag.

The [compatibility report](/guides/compatibility-report.html) records the
user-facing tools tested together. Library releases follow their own module
versions and support records.
