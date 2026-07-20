// Package site builds the PawnKit documentation site.
package site

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

const maxSourceSize int64 = 32 << 20

// Options controls a site build.
type Options struct {
	Config     string
	Output     string
	SourceRoot string
}

type config struct {
	Sources []source `json:"sources"`
}

type source struct {
	Name       string `json:"name"`
	Repository string `json:"repository"`
	Ref        string `json:"ref"`
	Path       string `json:"path"`
	SourcePath string `json:"source_path"`
	Kind       string `json:"kind"`
}

type entry struct {
	Title      string `json:"title"`
	URL        string `json:"url"`
	Kind       string `json:"kind"`
	Source     string `json:"source"`
	Version    string `json:"version"`
	Repository string `json:"repository"`
	RawURL     string `json:"raw_url,omitempty"`
	Summary    string `json:"summary,omitempty"`
}

type pageData struct {
	Title       string
	Description string
	Content     template.HTML
	Entry       entry
	Entries     []entry
}

var markdown = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
	goldmark.WithParserOptions(parser.WithAutoHeadingID()),
)

// Build creates a complete site in the output directory.
func Build(options Options) error {
	if options.Config == "" || options.Output == "" || options.SourceRoot == "" {
		return errors.New("config, output, and source root are required")
	}
	if err := validateOutput(options.Output); err != nil {
		return err
	}

	cfg, err := loadConfig(options.Config)
	if err != nil {
		return err
	}

	staging, err := os.MkdirTemp(filepath.Dir(options.Output), ".pawnkit-site-*")
	if err != nil {
		return err
	}
	defer func() { _ = os.RemoveAll(staging) }()

	if err := os.MkdirAll(filepath.Join(staging, "assets"), 0o750); err != nil {
		return err
	}
	if err := writeFile(filepath.Join(staging, "assets", "site.css"), []byte(styles)); err != nil {
		return err
	}
	if err := writeFile(filepath.Join(staging, "assets", "search.js"), []byte(searchScript)); err != nil {
		return err
	}
	if err := writeFile(filepath.Join(staging, "CNAME"), []byte("pawnkit.dev\n")); err != nil {
		return err
	}
	if err := writeFile(filepath.Join(staging, ".nojekyll"), nil); err != nil {
		return err
	}

	entries, err := buildSources(staging, options.SourceRoot, cfg.Sources)
	if err != nil {
		return err
	}
	guides, err := buildGuides(staging)
	if err != nil {
		return err
	}
	entries = append(entries, guides...)
	sortEntries(entries)

	if err := writeIndex(staging, entries); err != nil {
		return err
	}
	if err := writeSearch(staging, entries); err != nil {
		return err
	}
	if err := validateInternalLinks(staging); err != nil {
		return err
	}

	if err := os.RemoveAll(options.Output); err != nil {
		return err
	}
	return os.Rename(staging, options.Output)
}

func loadConfig(name string) (config, error) {
	data, err := os.ReadFile(name) //nolint:gosec // The caller supplies the configuration path.
	if err != nil {
		return config{}, err
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	var cfg config
	if err := decoder.Decode(&cfg); err != nil {
		return config{}, fmt.Errorf("read sources: %w", err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return config{}, errors.New("read sources: trailing JSON")
	}
	if len(cfg.Sources) == 0 {
		return config{}, errors.New("no documentation sources configured")
	}
	seen := make(map[string]bool)
	for _, item := range cfg.Sources {
		if item.Name == "" || item.Repository == "" || item.Ref == "" || item.Path == "" || item.SourcePath == "" || item.Kind == "" {
			return config{}, errors.New("each source needs a name, repository, ref, path, source path, and kind")
		}
		if seen[item.Name] {
			return config{}, fmt.Errorf("duplicate source %q", item.Name)
		}
		if !safeRelativePath(item.Path) || !safeRelativePath(item.SourcePath) || strings.ContainsAny(item.Kind, "/\\") || strings.ContainsAny(item.Ref, "/\\") {
			return config{}, fmt.Errorf("source %q has an unsafe path, kind, or ref", item.Name)
		}
		repository, err := url.Parse(item.Repository)
		if err != nil || repository.Scheme != "https" || repository.Host != "github.com" {
			return config{}, fmt.Errorf("source %q needs an https GitHub repository URL", item.Name)
		}
		seen[item.Name] = true
	}
	return cfg, nil
}

func safeRelativePath(name string) bool {
	if name == "" || strings.Contains(name, "\\") || path.IsAbs(name) {
		return false
	}
	clean := path.Clean(name)
	return clean == name && clean != ".." && !strings.HasPrefix(clean, "../")
}

func buildSources(output, sourceRoot string, sources []source) ([]entry, error) {
	var entries []entry
	for _, item := range sources {
		root := filepath.Join(sourceRoot, filepath.FromSlash(item.Path))
		info, err := os.Lstat(root)
		if err != nil {
			return nil, fmt.Errorf("source %s: %w", item.Name, err)
		}
		if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
			return nil, fmt.Errorf("source %s is not a directory", item.Name)
		}

		err = filepath.WalkDir(root, func(name string, file fs.DirEntry, walkErr error) error {
			if walkErr != nil || file.IsDir() {
				return walkErr
			}
			ext := strings.ToLower(filepath.Ext(name))
			if ext != ".json" && ext != ".md" {
				return nil
			}
			info, err := file.Info()
			if err != nil {
				return err
			}
			if !info.Mode().IsRegular() || info.Size() > maxSourceSize {
				return fmt.Errorf("unsupported source file %s", name)
			}
			relative, err := filepath.Rel(root, name)
			if err != nil {
				return err
			}
			relative = filepath.ToSlash(relative)
			data, err := os.ReadFile(name) //nolint:gosec // WalkDir selected a file under the configured source.
			if err != nil {
				return err
			}
			if ext == ".json" && !json.Valid(data) {
				return fmt.Errorf("invalid JSON in %s", name)
			}

			rawPath := path.Join("raw", item.Kind, item.Ref, relative)
			if err := writeFile(filepath.Join(output, filepath.FromSlash(rawPath)), data); err != nil {
				return err
			}
			latestPath := path.Join("raw", item.Kind, "latest", relative)
			if err := writeFile(filepath.Join(output, filepath.FromSlash(latestPath)), data); err != nil {
				return err
			}
			publicRawURL := "/" + rawPath
			if item.Kind == "schema" {
				schemaURL, err := writeSchemaURL(output, data)
				if err != nil {
					return fmt.Errorf("schema %s: %w", name, err)
				}
				publicRawURL = schemaURL
			}

			title := titleFromName(relative)
			if ext == ".md" {
				title = titleFromMarkdown(data, relative)
			}
			url := "/" + rawPath
			summary := firstParagraph(data)
			if ext == ".json" {
				summary = fmt.Sprintf("%s data from %s.", titleFromName(item.Kind), item.Name)
			}
			if item.Kind == "schema" {
				url = publicRawURL
			}
			if ext == ".md" {
				url = "/reference/" + item.Kind + "/" + strings.TrimSuffix(relative, ext) + ".html"
			}
			entry := entry{
				Title: title, URL: url, RawURL: publicRawURL, Kind: item.Kind,
				Source: item.Name, Version: item.Ref, Repository: sourceFileURL(item, relative),
				Summary: summary,
			}
			if ext == ".md" {
				body, err := renderMarkdown(rewriteMarkdownLinks(data, item, relative))
				if err != nil {
					return err
				}
				page := pageData{Title: title, Description: entry.Summary, Content: body, Entry: entry}
				pageName := filepath.FromSlash(strings.TrimSuffix(relative, ext) + ".html")
				for _, prefix := range []string{"", item.Ref, "latest"} {
					target := filepath.Join(output, "reference", item.Kind, prefix, pageName)
					if err := renderPage(target, page); err != nil {
						return err
					}
				}
			}
			entries = append(entries, entry)
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return entries, nil
}

func sourceFileURL(item source, relative string) string {
	return strings.TrimSuffix(item.Repository, "/") + "/blob/" + item.Ref + "/" + path.Join(item.SourcePath, relative)
}

func rewriteMarkdownLinks(data []byte, item source, relative string) []byte {
	return markdownLinkPattern.ReplaceAllFunc(data, func(match []byte) []byte {
		parts := markdownLinkPattern.FindSubmatch(match)
		if len(parts) != 2 {
			return match
		}
		target := string(parts[1])
		if target == "" || strings.HasPrefix(target, "#") {
			return match
		}
		parsed, err := url.Parse(target)
		if err != nil || parsed.IsAbs() {
			return match
		}
		clean := path.Clean(path.Join(path.Dir(relative), parsed.Path))
		if strings.HasSuffix(parsed.Path, ".md") && clean != ".." && !strings.HasPrefix(clean, "../") {
			parsed.Path = strings.TrimSuffix(parsed.Path, ".md") + ".html"
			return bytes.Replace(match, parts[1], []byte(parsed.String()), 1)
		}
		repositoryPath := path.Clean(path.Join(item.SourcePath, path.Dir(relative), parsed.Path))
		if repositoryPath == ".." || strings.HasPrefix(repositoryPath, "../") {
			return match
		}
		resolved := strings.TrimSuffix(item.Repository, "/") + "/blob/" + item.Ref + "/" + repositoryPath
		return bytes.Replace(match, parts[1], []byte(resolved), 1)
	})
}

func validateOutput(output string) error {
	absolute, err := filepath.Abs(output)
	if err != nil {
		return err
	}
	working, err := os.Getwd()
	if err != nil {
		return err
	}
	relative, err := filepath.Rel(absolute, working)
	if err != nil {
		return err
	}
	if relative == "." || (relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator))) {
		return errors.New("output cannot contain the working directory")
	}
	return nil
}

func writeSchemaURL(output string, data []byte) (string, error) {
	var document struct {
		ID string `json:"$id"`
	}
	if err := json.Unmarshal(data, &document); err != nil {
		return "", err
	}
	parsed, err := url.Parse(document.ID)
	if err != nil {
		return "", errors.New("$id must be an https://schemas.pawnkit.dev URL")
	}
	schemaPath := strings.TrimPrefix(parsed.Path, "/")
	if parsed.Scheme != "https" || parsed.Host != "schemas.pawnkit.dev" || parsed.RawQuery != "" || parsed.Fragment != "" || !safeRelativePath(schemaPath) {
		return "", errors.New("$id must be an https://schemas.pawnkit.dev URL")
	}
	if err := writeFile(filepath.Join(output, filepath.FromSlash(schemaPath)), data); err != nil {
		return "", err
	}
	return document.ID, nil
}

func buildGuides(output string) ([]entry, error) {
	files, err := filepath.Glob(filepath.Join("content", "*.md"))
	if err != nil {
		return nil, err
	}
	var entries []entry
	for _, name := range files {
		data, err := os.ReadFile(name) //nolint:gosec // Glob limits files to authored guides.
		if err != nil {
			return nil, err
		}
		slug := strings.TrimSuffix(filepath.Base(name), ".md")
		entry := entry{Title: titleFromMarkdown(data, slug), URL: "/guides/" + slug + ".html", Kind: "guide", Source: "pawnkit.dev", Version: "current", Summary: firstParagraph(data)}
		body, err := renderMarkdown(data)
		if err != nil {
			return nil, err
		}
		if err := renderPage(filepath.Join(output, "guides", slug+".html"), pageData{Title: entry.Title, Description: entry.Summary, Content: body, Entry: entry}); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func renderMarkdown(data []byte) (template.HTML, error) {
	var output bytes.Buffer
	if err := markdown.Convert(data, &output); err != nil {
		return "", err
	}
	return template.HTML(output.String()), nil //nolint:gosec // Goldmark escapes raw HTML by default.
}

func renderPage(name string, data pageData) error {
	if err := os.MkdirAll(filepath.Dir(name), 0o750); err != nil {
		return err
	}
	var output bytes.Buffer
	if err := pageTemplate.Execute(&output, data); err != nil {
		return err
	}
	return writeFile(name, output.Bytes())
}

func writeIndex(output string, entries []entry) error {
	grouped := make(map[string][]entry)
	for _, item := range entries {
		grouped[item.Kind] = append(grouped[item.Kind], item)
	}
	var content strings.Builder
	content.WriteString("<section class=\"hero\"><p class=\"eyebrow\">Pawn tooling, in one place</p><h1>Build and maintain Pawn projects without guesswork.</h1><p>Start with a practical guide, or search the standards and generated references used by PawnKit tools.</p></section>")
	content.WriteString("<section><h2>Guides</h2><div class=\"cards\">")
	for _, item := range grouped["guide"] {
		fmt.Fprintf(&content, "<article><h3><a href=\"%s\">%s</a></h3><p>%s</p></article>", template.HTMLEscapeString(item.URL), template.HTMLEscapeString(item.Title), template.HTMLEscapeString(item.Summary))
	}
	content.WriteString("</div></section><section><h2>Reference</h2><p>Browse RFCs, schemas, API data, and lint rules from their owning repositories.</p><ul class=\"reference-links\">")
	for _, kind := range []string{"rfc", "schema", "api", "rule"} {
		if len(grouped[kind]) > 0 {
			fmt.Fprintf(&content, "<li><a href=\"/search.html?kind=%s\">%s <span>%d</span></a></li>", kind, template.HTMLEscapeString(strings.ToUpper(kind)), len(grouped[kind]))
		}
	}
	content.WriteString("</ul></section>")
	if err := renderPage(filepath.Join(output, "index.html"), pageData{Title: "PawnKit documentation", Description: "Guides and generated reference material for PawnKit.", Content: template.HTML(content.String())}); err != nil { //nolint:gosec // Content is escaped above.
		return err
	}
	return renderPage(filepath.Join(output, "search.html"), pageData{Title: "Search", Description: "Search PawnKit guides and reference material.", Content: template.HTML(searchMarkup)}) //nolint:gosec // Static markup is trusted.
}

func writeSearch(output string, entries []entry) error {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return writeFile(filepath.Join(output, "search.json"), append(data, '\n'))
}

func validateInternalLinks(root string) error {
	rootHandle, err := os.OpenRoot(root)
	if err != nil {
		return err
	}
	defer func() { _ = rootHandle.Close() }()

	return filepath.WalkDir(root, func(name string, file fs.DirEntry, walkErr error) error {
		if walkErr != nil || file.IsDir() || filepath.Ext(name) != ".html" {
			return walkErr
		}
		data, err := os.ReadFile(name) //nolint:gosec // WalkDir selected generated HTML.
		if err != nil {
			return err
		}
		for _, match := range internalLinkPattern.FindAllSubmatch(data, -1) {
			target := strings.SplitN(string(match[1]), "#", 2)[0]
			if target == "" || strings.HasPrefix(target, "#") || strings.Contains(target, "://") || strings.HasPrefix(target, "mailto:") {
				continue
			}
			target = strings.SplitN(target, "?", 2)[0]
			if target == "/" {
				target = "/index.html"
			}
			if strings.HasSuffix(target, "/") {
				target += "index.html"
			}
			if !strings.HasPrefix(target, "/") {
				relative, err := filepath.Rel(root, filepath.Dir(name))
				if err != nil {
					return err
				}
				target = filepath.ToSlash(filepath.Join(relative, filepath.FromSlash(target)))
			}
			if _, err := rootHandle.Stat(filepath.FromSlash(strings.TrimPrefix(target, "/"))); err != nil {
				return fmt.Errorf("broken link in %s: %s", name, target)
			}
		}
		return nil
	})
}

func writeFile(name string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(name), 0o750); err != nil {
		return err
	}
	return os.WriteFile(name, data, 0o644) //nolint:gosec // Site files are public.
}

func sortEntries(entries []entry) {
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Kind != entries[j].Kind {
			return entries[i].Kind < entries[j].Kind
		}
		return entries[i].Title < entries[j].Title
	})
}

func titleFromMarkdown(data []byte, fallback string) string {
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "# ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "# "))
		}
	}
	return titleFromName(fallback)
}

func titleFromName(name string) string {
	name = strings.TrimSuffix(filepath.Base(name), filepath.Ext(name))
	name = strings.ReplaceAll(name, "-", " ")
	name = strings.ReplaceAll(name, "_", " ")
	runes := []rune(name)
	if len(runes) > 0 {
		runes[0] = unicode.ToUpper(runes[0])
	}
	return string(runes)
}

func firstParagraph(data []byte) string {
	lines := strings.Split(string(data), "\n")
	var paragraph []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			if len(paragraph) > 0 {
				break
			}
			continue
		}
		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "```") {
			continue
		}
		paragraph = append(paragraph, line)
	}
	text := strings.Join(paragraph, " ")
	text = strings.NewReplacer("`", "", "**", "", "*", "").Replace(text)
	if len(text) > 180 {
		text = strings.TrimSpace(text[:177]) + "..."
	}
	return text
}
