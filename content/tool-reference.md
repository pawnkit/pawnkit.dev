# Choose a tool

Start with the problem you are trying to solve. Most projects only need a few of
these tools.

| Tool | Use it for |
| --- | --- |
| `pawn` | Project checks, environment diagnostics, and shared workflows |
| `pawnfmt` | Stable source formatting |
| `pawnlint` | Semantic and style diagnostics |
| `pawntest` | Pawn unit tests, coverage, and profiling |
| `pawnlsp` | Editor diagnostics, navigation, and completion |
| `pawnmigrate` | Reviewable source and API migrations |
| `pawndoc` | API documentation generated from Pawn source |
| `pawndebug` | AMX debugging through the Debug Adapter Protocol |
| `pawnserver` | Server bundle inspection, installation, and rollback |

Run a binary with `--help` for its exact command and flag list. Released command
help is authoritative; this site focuses on workflows that cross tool boundaries.

Library repositories such as `pawnkit-core`, `pawn-parser`, `pawn-analysis`, and
`pawn-project` are intended for Go programs that extend the ecosystem.
