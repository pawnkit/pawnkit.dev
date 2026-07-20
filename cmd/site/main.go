package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/pawnkit/pawnkit.dev/internal/site"
)

func main() {
	defaultSourceRoot := os.Getenv("PAWNKIT_SOURCE_ROOT")
	if defaultSourceRoot == "" {
		defaultSourceRoot = ".."
	}
	config := flag.String("config", "sources.json", "source configuration")
	output := flag.String("output", "dist", "output directory")
	root := flag.String("source-root", defaultSourceRoot, "directory containing source repositories")
	flag.Parse()

	if err := site.Build(site.Options{Config: *config, Output: *output, SourceRoot: *root}); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
