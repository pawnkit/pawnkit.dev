# CI and editors

Use the same checks locally and in CI. This keeps editor feedback from drifting
away from pull-request results.

## GitHub Actions

The PawnKit check action installs the released CLI and runs the project workflow:

```yaml
name: Pawn

on:
  push:
  pull_request:

permissions:
  contents: read

jobs:
  check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: pawnkit/pawn-actions/check@v1
```

Pin action commits when your project requires a locked supply chain.

## VS Code

Install the Pawn extension, open the project root, and run `Pawn: Doctor` from
the command palette. The extension finds `pawnlsp` and the other PawnKit binaries
on `PATH`.

Keep project settings in `pawn.json`. Workspace-only editor settings are useful
for personal preferences, but shared compiler paths and include directories
belong in the project manifest.
