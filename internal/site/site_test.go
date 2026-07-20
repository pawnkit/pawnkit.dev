package site

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSafeRelativePath(t *testing.T) {
	tests := map[string]bool{
		"pawnkit-spec/rfcs": true,
		"schemas":           true,
		"../schemas":        false,
		"/schemas":          false,
		`repo\schemas`:      false,
		"repo/../schemas":   false,
	}
	for name, want := range tests {
		if got := safeRelativePath(name); got != want {
			t.Errorf("safeRelativePath(%q) = %v, want %v", name, got, want)
		}
	}
}

func TestLoadConfigRejectsTrailingJSON(t *testing.T) {
	name := filepath.Join(t.TempDir(), "sources.json")
	if err := os.WriteFile(name, []byte(`{"sources":[]} {}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := loadConfig(name); err == nil {
		t.Fatal("trailing JSON was accepted")
	}
}

func TestBuildRejectsWorkingDirectoryAsOutput(t *testing.T) {
	if err := validateOutput("."); err == nil {
		t.Fatal("working directory was accepted as output")
	}
}

func TestRewriteMarkdownLinks(t *testing.T) {
	item := source{Repository: "https://github.com/pawnkit/pawnlint", Ref: "v1.0.2", SourcePath: "docs/rules"}
	input := []byte("[rule](other.md) [guide](../guide.md) [web](https://example.com)")
	want := "[rule](other.html) [guide](https://github.com/pawnkit/pawnlint/blob/v1.0.2/docs/guide.md) [web](https://example.com)"
	if got := string(rewriteMarkdownLinks(input, item, "index.md")); got != want {
		t.Fatalf("links = %q, want %q", got, want)
	}
}

func TestWriteSchemaURLRejectsTraversal(t *testing.T) {
	data := []byte(`{"$id":"https://schemas.pawnkit.dev/../outside.json"}`)
	if _, err := writeSchemaURL(t.TempDir(), data); err == nil {
		t.Fatal("traversing schema ID was accepted")
	}
}
