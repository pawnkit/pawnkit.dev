package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadCheckoutsDeduplicatesRepository(t *testing.T) {
	config := `{"sources":[
		{"repository":"https://github.com/pawnkit/pawnkit-spec","ref":"v1.2.3","path":"pawnkit-spec/schemas"},
		{"repository":"https://github.com/pawnkit/pawnkit-spec","ref":"v1.2.3","path":"pawnkit-spec/rfcs"}
	]}`
	path := filepath.Join(t.TempDir(), "sources.json")
	if err := os.WriteFile(path, []byte(config), 0o600); err != nil {
		t.Fatal(err)
	}

	checkouts, err := loadCheckouts(path, "sources")
	if err != nil {
		t.Fatal(err)
	}
	if len(checkouts) != 1 || checkouts[0].ref != "v1.2.3" {
		t.Fatalf("checkouts = %#v", checkouts)
	}
}

func TestLoadCheckoutsRejectsConflictingPins(t *testing.T) {
	config := `{"sources":[
		{"repository":"https://github.com/pawnkit/pawnkit-spec","ref":"v1.2.3","path":"pawnkit-spec/schemas"},
		{"repository":"https://github.com/pawnkit/pawnkit-spec","ref":"v1.2.4","path":"pawnkit-spec/rfcs"}
	]}`
	path := filepath.Join(t.TempDir(), "sources.json")
	if err := os.WriteFile(path, []byte(config), 0o600); err != nil {
		t.Fatal(err)
	}

	if _, err := loadCheckouts(path, "sources"); err == nil {
		t.Fatal("conflicting pins were accepted")
	}
}
