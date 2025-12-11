// Package main provides txt2llm, a tool for concatenating files with markers for LLM input.
package main

import (
	"fmt"
	"os"

	"github.com/matthewchivers/txt2llm/pkg/cli"
	"github.com/matthewchivers/txt2llm/pkg/output"
	"github.com/matthewchivers/txt2llm/pkg/resolve"
)

func main() {
	cfg := cli.Parse()
	patterns := cli.Patterns()
	files, err := resolve.Files(patterns, cfg.Recursive)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	outPaths := output.Paths(files, cfg.Relative)
	output.Header(cfg.MarkerPrefix, cfg.MarkerSuffix)
	output.Markers(files, outPaths, cfg.MarkerPrefix, cfg.MarkerSuffix)
}
