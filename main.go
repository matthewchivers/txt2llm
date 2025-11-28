package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"
)

func main() {
	// Flags
	var recursive bool
	pflag.BoolVar(&recursive, "recursive", false, "Process directories recursively")
	var relative bool
	pflag.BoolVar(&relative, "relative", false, "Use paths relative to current directory in output")
	var markerPrefix string
	var markerSuffix string
	pflag.StringVar(&markerPrefix, "marker-prefix", "<<<", "Prefix for start/end marker lines")
	pflag.StringVar(&markerSuffix, "marker-suffix", ">>>", "Suffix for start/end marker lines")

	// Allow flags to be interspersed with positional args
	pflag.CommandLine.SetInterspersed(true)
	pflag.Parse()

	// Use space-separated positional args as patterns (files, dirs, globs)
	patterns := append([]string{}, pflag.Args()...)

	// Resolve patterns into files
	files := resolveFiles(patterns, recursive)
	if len(files) == 0 {
		fmt.Println("No input files matched. Provide patterns like 'file.txt' 'dir' '*.md'.")
		os.Exit(1)
	}

	// Optionally convert paths to relative for output only
	outPaths := files
	if relative {
		cwd, _ := os.Getwd()
		rels := make([]string, 0, len(files))
		for _, p := range files {
			if r, err := filepath.Rel(cwd, p); err == nil {
				rels = append(rels, r)
			} else {
				rels = append(rels, p)
			}
		}
		outPaths = rels
	}

	// Print a simple key at the top explaining markers and filename
	printKey(markerPrefix, markerSuffix)

	// Output files with markers
	outputFilesWithMarkers(files, outPaths, markerPrefix, markerSuffix)
}

func outputFilesWithMarkers(files []string, outPaths []string, markerPrefix, markerSuffix string) {
	for i, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", f, err)
			continue
		}
		fmt.Printf("%sSTART:%s%s\n", markerPrefix, outPaths[i], markerSuffix)
		os.Stdout.Write(data)
		if len(data) > 0 && data[len(data)-1] != '\n' {
			fmt.Println()
		}
		fmt.Printf("%sEND:%s%s\n\n", markerPrefix, outPaths[i], markerSuffix)
	}
}

// printKey prints a simple key explaining how to read markers and filenames.
func printKey(markerPrefix, markerSuffix string) {
	fmt.Printf("Each section below represents text output from one file.\n")
	fmt.Printf("Delimiters: %sSTART:{filename}%s ... %sEND:{filename}%s\n\n", markerPrefix, markerSuffix, markerPrefix, markerSuffix)
}

// resolveFiles takes a list of patterns (files, directories, globs) and returns a
// deduplicated list of file paths. If recursive is true, directories are processed recursively.
func resolveFiles(patterns []string, recursive bool) []string {
	seen := map[string]struct{}{}
	out := []string{}
	add := func(path string) {
		if path == "" {
			return
		}
		abs, err := filepath.Abs(path)
		if err == nil {
			path = abs
		}
		if _, ok := seen[path]; !ok {
			seen[path] = struct{}{}
			out = append(out, path)
		}
	}

	for _, pat := range patterns {
		if pat == "" {
			continue
		}
		// If pattern points to an existing path, handle according to type
		if info, err := os.Stat(pat); err == nil {
			if info.Mode().IsRegular() {
				add(pat)
				continue
			}
			if info.IsDir() {
				if recursive {
					filepath.WalkDir(pat, func(path string, d os.DirEntry, err error) error {
						if err != nil {
							return nil
						}
						if d.Type().IsRegular() {
							add(path)
						}
						return nil
					})
				} else {
					entries, _ := os.ReadDir(pat)
					for _, e := range entries {
						if e.Type().IsRegular() {
							add(filepath.Join(pat, e.Name()))
						}
					}
				}
				continue
			}
		}
		// Otherwise treat as glob pattern
		matches, _ := filepath.Glob(pat)
		for _, m := range matches {
			if info, err := os.Stat(m); err == nil && info.Mode().IsRegular() {
				add(m)
			}
		}
	}

	return out
}
