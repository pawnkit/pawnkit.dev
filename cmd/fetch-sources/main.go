package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

type config struct {
	Sources []source `json:"sources"`
}

type source struct {
	Repository string `json:"repository"`
	Ref        string `json:"ref"`
	Path       string `json:"path"`
}

type checkout struct {
	repository string
	ref        string
	directory  string
}

func main() {
	configPath := flag.String("config", "sources.json", "source configuration")
	output := flag.String("output", "sources", "checkout directory")
	flag.Parse()

	checkouts, err := loadCheckouts(*configPath, *output)
	if err != nil {
		fail(err)
	}
	for _, item := range checkouts {
		if err := clone(item); err != nil {
			fail(err)
		}
	}
}

func loadCheckouts(configPath, output string) ([]checkout, error) {
	data, err := os.ReadFile(configPath) //nolint:gosec // The caller selects the config.
	if err != nil {
		return nil, err
	}
	var cfg config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if output == "" || filepath.IsAbs(output) || filepath.Clean(output) != output {
		return nil, errors.New("output must be a clean relative path")
	}

	seen := make(map[string]checkout)
	for _, item := range cfg.Sources {
		parts := strings.Split(item.Path, "/")
		if len(parts) < 2 || parts[0] == "" || path.Clean(item.Path) != item.Path {
			return nil, fmt.Errorf("invalid source path %q", item.Path)
		}
		if !strings.HasPrefix(item.Repository, "https://github.com/pawnkit/") || item.Ref == "" ||
			strings.HasPrefix(item.Ref, "-") || strings.ContainsAny(item.Ref, "/\\") {
			return nil, fmt.Errorf("invalid source repository or ref for %q", item.Path)
		}
		next := checkout{repository: item.Repository, ref: item.Ref, directory: filepath.Join(output, parts[0])}
		if current, ok := seen[parts[0]]; ok {
			if current.repository != next.repository || current.ref != next.ref {
				return nil, fmt.Errorf("conflicting pins for %s", parts[0])
			}
			continue
		}
		seen[parts[0]] = next
	}

	out := make([]checkout, 0, len(seen))
	for _, item := range seen {
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].directory < out[j].directory })

	return out, nil
}

func clone(item checkout) error {
	if _, err := os.Lstat(item.directory); err == nil {
		return fmt.Errorf("source directory already exists: %s", item.directory)
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(item.directory), 0o750); err != nil {
		return err
	}
	command := exec.Command("git", "clone", "--quiet", "--filter=blob:none", "--depth", "1", "--branch", item.ref, "--", item.repository, item.directory) //nolint:gosec // Config accepts only pinned PawnKit repositories.
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	if err := command.Run(); err != nil {
		return fmt.Errorf("clone %s at %s: %w", item.repository, item.ref, err)
	}
	return nil
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
