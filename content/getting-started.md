# Getting started

PawnKit tools work with existing Pawn projects. You do not need to move files or
adopt every tool at once.

## Install the CLI

Download `pawn` from the
[pawnkit-cli releases](https://github.com/pawnkit/pawnkit-cli/releases), then put
it on your `PATH`.

Check the project from its root directory:

```sh
pawn doctor
pawn check
```

`pawn doctor` reports missing tools and configuration problems. `pawn check`
runs the configured formatter, linter, tests, and compiler checks.

## Describe the project

Create `pawn.json` in the project root. Start with the target and entry point:

```json
{
  "$schema": "https://schemas.pawnkit.dev/pawn-project/v1/schema.json",
  "target": "openmp",
  "entry": "gamemodes/main.pwn"
}
```

Add include directories and tool settings only when the project needs them. The
[project manifest RFC](/reference/rfc/0002-project-manifest.html) lists every
field.

## Add tools gradually

Formatting is a low-risk first step:

```sh
pawn fmt --check
pawn fmt
```

Then add linting and tests. Keep the commands in CI once the existing codebase
passes them consistently.
