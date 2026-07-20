package site

import (
	_ "embed"
	"html/template"
	"regexp"
)

var (
	internalLinkPattern = regexp.MustCompile(`(?:href|src)=["']([^"']+)`)
	markdownLinkPattern = regexp.MustCompile(`\]\(([^)[:space:]]+)\)`)
)

//go:embed static/page.html
var pageHTML string

//go:embed static/site.css
var styles string

//go:embed static/search.js
var searchScript string

var pageTemplate = template.Must(template.New("page").Parse(pageHTML))

const searchMarkup = `<h1>Search</h1><p>Search guides, RFCs, schemas, API data, and lint rules.</p><form role="search"><label for="search-query">Search documentation</label><input id="search-query" name="q" type="search" autocomplete="off"><label for="search-kind">Type</label><select id="search-kind" name="kind"><option value="">Everything</option><option value="guide">Guides</option><option value="rfc">RFCs</option><option value="schema">Schemas</option><option value="api">API data</option><option value="rule">Lint rules</option></select></form><p id="search-status" aria-live="polite"></p><ul id="search-results" class="search-results"></ul>`
