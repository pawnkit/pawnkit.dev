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
