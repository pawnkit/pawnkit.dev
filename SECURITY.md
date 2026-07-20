# Security policy

Report vulnerabilities through a private
[GitHub security advisory](https://github.com/pawnkit/pawnkit.dev/security/advisories/new).
Do not open a public issue before a fix is available.

The build reads pinned local checkouts and does not scrape sites at runtime.
Configured source files are limited to 32 MiB, Markdown renders with raw HTML
disabled, and the deployed site has no server-side code or credentials.
