# Deployment

The `Deploy` workflow builds the site from pinned source checkouts and publishes
`dist` to GitHub Pages whenever `main` changes.

Before the first deployment:

1. Make `pawnkit-spec`, `pawn-api`, and `pawnlint` publicly readable.
2. Set the repository's Pages source to GitHub Actions.
3. Point `pawnkit.dev` at GitHub Pages and enable HTTPS.
4. Route `schemas.pawnkit.dev` to the same built files without changing the URL.

The last route matters because schema `$id` values use the schema subdomain. A
DNS record alone may not be enough if the hosting provider only accepts one
custom hostname. A small reverse proxy or an additional static-site mapping can
serve the same deployment for both names.

Pull requests run the complete build but do not deploy. Download the workflow
artifact or run `go run ./cmd/site` for a local preview.

When updating a source, change its ref in both `sources.json` and the checkout
steps. The build fails if expected files, schema IDs, or internal links are
missing.
