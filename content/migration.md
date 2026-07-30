# Migrate an existing project

Commit or stash current changes before running a migration. `pawnmigrate` checks
the working tree, but a clean commit is the easiest rollback.

Preview the plan first:

```sh
pawnmigrate
```

Read the proposed edits, then apply them:

```sh
pawnmigrate --apply
```

Run the formatter, linter, tests, and compiler after each migration group. Keep
API deprecations separate from mechanical syntax changes; they usually need more
review and may change runtime behaviour.

For an SA-MP to open.mp move, update the target profile and dependencies before
rewriting APIs. This gives the analysis tools the right includes and API metadata
for the rest of the migration.

PawnKit keeps reading sampctl-compatible projects so adoption does not require
an immediate manifest rewrite. That compatibility is a migration path: the
target workflow uses `pawn` for dependencies, toolchains, builds, tests, and
runtime commands without requiring sampctl.

`pawnmigrate` also handles reviewed source migrations, including moves between
Pawn libraries when a supported rule exists. These rules may change project
dependencies and source together. They leave ambiguous code unchanged and
report work that still needs a maintainer.
