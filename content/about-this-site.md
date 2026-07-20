# About this site

PawnKit is maintained by people building tools for the Pawn, open.mp, and SA-MP
communities in their spare time. This site puts the shared guides and reference
material in one place.

## Which page is authoritative?

Guides are maintained here. Schemas, RFCs, API data, and lint rules come from the
repository named at the top of each reference page. Fix a generated fact in that
repository so every consumer receives the same correction.

Source revisions are pinned in `sources.json`. Reference pages and raw files are
published under both a revision and `latest`. Use a versioned URL when a build
must remain reproducible.

Schema `$id` URLs use `schemas.pawnkit.dev`. They are stable contracts and do not
follow draft changes until the schema version changes.

## Releases and compatibility

The site itself follows semantic versioning for its build pipeline. Guide edits
and newly indexed material are not breaking changes. Moving or removing a stable
raw URL is breaking and requires a redirect or a new major version.

Pre-1.0 source projects may still change their formats. Their pinned revision is
shown on every generated page so you can check exactly what the site published.
