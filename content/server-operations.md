# Run a packaged server

`pawnserver` installs versioned open.mp and SA-MP server bundles. A bundle records
its binaries, checksums, entry points, configuration, and persistent paths.

Inspect an archive before installing it:

```sh
pawnserver inspect release.tar.gz
pawnserver verify release.tar.gz
pawnserver install release.tar.gz /srv/my-server
```

Installation prints a plan unless `--apply` is supplied. Review the destination
and version before applying it:

```sh
pawnserver install --apply release.tar.gz /srv/my-server
```

Updates keep the manifest's persistent paths and retain the previous installation
for rollback:

```sh
pawnserver update --apply release.tar.gz /srv/my-server
pawnserver rollback /srv/my-server
```

Keep secrets and mutable databases outside the bundle. Run `pawnserver doctor`
after deployment to verify checksums, entry points, and configuration.

See the [server bundle RFC](/reference/rfc/0006-server-bundle.html) for the file
format and platform names.
