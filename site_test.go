package site_test

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

type searchEntry struct {
	Title      string `json:"title"`
	URL        string `json:"url"`
	Kind       string `json:"kind"`
	Version    string `json:"version"`
	Repository string `json:"repository"`
	RawURL     string `json:"raw_url"`
}

func TestBuildAndSearchIndex(t *testing.T) {
	command := exec.Command("go", "run", "./cmd/site")
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("build: %v: %s", err, output)
	}

	data, err := os.ReadFile("dist/search.json")
	if err != nil {
		t.Fatal(err)
	}
	var entries []searchEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		t.Fatal(err)
	}
	if len(entries) < 10 {
		t.Fatalf("search index has only %d entries", len(entries))
	}

	var foundRule bool
	for _, entry := range entries {
		if entry.URL == "/reference/rule/missing-include.html" {
			foundRule = true
			if entry.Kind != "rule" || entry.Version != "v1.0.2" || entry.Repository == "" || entry.RawURL == "" {
				t.Fatalf("rule provenance = %#v", entry)
			}
		}
	}
	if !foundRule {
		t.Fatal("missing generated lint rule")
	}

	assertContains(t, "dist/index.html", `class="skip-link"`)
	assertContains(t, "dist/index.html", `href="/guides/getting-started.html"`)
	assertContains(t, "dist/search.html", `aria-live="polite"`)
	assertContains(t, "dist/reference/rule/missing-include.html", "Raw file")

	if _, err := os.Stat(filepath.Join("dist", "raw", "rule", "v1.0.2", "missing-include.md")); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join("dist", "raw", "rule", "latest", "missing-include.md")); err != nil {
		t.Fatal(err)
	}
}

func TestBuildIsDeterministic(t *testing.T) {
	command := exec.Command("go", "run", "./cmd/site", "-output", "dist-one")
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("first build: %v: %s", err, output)
	}
	t.Cleanup(func() { _ = os.RemoveAll("dist-one") })

	command = exec.Command("go", "run", "./cmd/site", "-output", "dist-two")
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("second build: %v: %s", err, output)
	}
	t.Cleanup(func() { _ = os.RemoveAll("dist-two") })

	for _, name := range []string{"index.html", "search.html", "search.json", "assets/site.css", "assets/search.js"} {
		one, err := os.ReadFile(filepath.Join("dist-one", name))
		if err != nil {
			t.Fatal(err)
		}
		two, err := os.ReadFile(filepath.Join("dist-two", name))
		if err != nil {
			t.Fatal(err)
		}
		if string(one) != string(two) {
			t.Fatalf("%s changed between builds", name)
		}
	}
}

func assertContains(t *testing.T, name, want string) {
	t.Helper()
	data, err := os.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), want) {
		t.Fatalf("%s does not contain %q", name, want)
	}
}
