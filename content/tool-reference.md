# Choose a PawnKit product

Most Pawn developers need the editor extension, the `pawn` command, or both.
The smaller tools are included behind those two entry points.

## Writing Pawn in an editor

Install the PawnKit extension for VS Code. It manages a tested set of the
language server, formatter, linter, and test runner.

Use the extension when you want diagnostics, completion, navigation,
formatting, tests, and project commands in the editor.

## Working in a terminal or CI

Install the `pawn` command. It runs the same project workflow locally and in
GitHub Actions:

```sh
pawn doctor
pawn check
pawn build
pawn test
pawn run
```

Use `pawnmigrate` for source migrations and `pawndoc` to generate API
documentation. These are separate commands because they change or generate
project files.

## Building Pawn tooling

Use the PawnKit Go libraries if you are writing an editor, build tool, analyser,
or runtime integration. Start with the [SDK guide](/guides/sdk.html); it maps
each job to its owning package.

## Running a project

Use `pawn run` for local open.mp development. It builds the project and starts a
verified runtime in an isolated session.

`pawnserver` handles versioned deployment bundles. See
[Run a packaged server](/guides/server-operations.html).

Run any command with `--help` for its current flags. The
[compatibility report](/guides/compatibility-report.html) lists the tool
versions tested together.
