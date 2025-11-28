package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"
)

func main() {
	// get configuration from flags
	cfg := parseFlags()
	patterns := append([]string{}, pflag.Args()...)

	// resolve input files from patterns
	files := resolveFiles(patterns, cfg.recursive)
	if len(files) == 0 {
		fmt.Println("No input files matched. Provide patterns like 'file.txt' 'dir' '*.md'.")
		os.Exit(1)
	}
	outPaths := getFilePaths(files, cfg.relative)

	// output files with markers
	printKey(cfg.markerPrefix, cfg.markerSuffix)
	outputFilesWithMarkers(files, outPaths, cfg.markerPrefix, cfg.markerSuffix)
}

// config holds command-line configuration options.
type config struct {
	recursive    bool   // process directories recursively
	relative     bool   // use relative paths in output
	markerPrefix string // prefix for start/end marker lines
	markerSuffix string // suffix for start/end marker lines
}

// parseFlags parses command-line flags and returns a populated config struct.
func parseFlags() config {
	var cfg config
	pflag.BoolVar(&cfg.recursive, "recursive", false, "Process directories recursively")
	pflag.BoolVar(&cfg.relative, "relative", false, "Use paths relative to current directory in output")
	pflag.StringVar(&cfg.markerPrefix, "marker-prefix", "<<<", "Prefix for start/end marker lines")
	pflag.StringVar(&cfg.markerSuffix, "marker-suffix", ">>>", "Suffix for start/end marker lines")
	pflag.CommandLine.SetInterspersed(true)
	pflag.Parse()
	return cfg
}

// getFilePaths returns paths for specified files, either absolute or relative to cwd.
func getFilePaths(files []string, relative bool) []string {
	if !relative {
		return files
	}
	cwd, _ := os.Getwd()
	rels := make([]string, 0, len(files))
	for _, p := range files {
		if r, err := filepath.Rel(cwd, p); err == nil {
			rels = append(rels, r)
		} else {
			rels = append(rels, p)
		}
	}
	return rels
}

// outputFilesWithMarkers outputs each file wrapped with start and end markers.
func outputFilesWithMarkers(files []string, outPaths []string, markerPrefix, markerSuffix string) {
	for i, f := range files {
		emitSection(f, outPaths[i], markerPrefix, markerSuffix)
	}
}

// emitSection reads the file at srcPath and writes it to stdout
// wrapped with start and end markers including outPath.
func emitSection(srcPath, outPath, markerPrefix, markerSuffix string) {
	data, err := os.ReadFile(srcPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", srcPath, err)
		return
	}
	fmt.Printf("%sSTART:%s%s\n", markerPrefix, outPath, markerSuffix)
	os.Stdout.Write(data)
	ensureTrailingNewline(data)
	fmt.Printf("%sEND:%s%s\n\n", markerPrefix, outPath, markerSuffix)
}

// ensureTrailingNewline prints a newline if data does not end with one.
func ensureTrailingNewline(data []byte) {
	if len(data) == 0 || data[len(data)-1] == '\n' {
		return
	}
	fmt.Println()
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
		if err != nil {
			abs = path
		}
		if _, ok := seen[abs]; !ok {
			seen[abs] = struct{}{}
			out = append(out, abs)
		}
	}

	for _, pat := range patterns {
		if pat == "" {
			continue
		}
		if info, err := os.Stat(pat); err == nil {
			if info.Mode().IsRegular() {
				add(pat)
				continue
			}
			if info.IsDir() {
				addDir(pat, recursive, add)
				continue
			}
		}
		addGlob(pat, add)
	}
	return out
}

// addDir adds regular files from the specified directory.
// If recursive is true, it walks the directory tree.
func addDir(dir string, recursive bool, add func(string)) {
	if recursive {
		filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if d.Type().IsRegular() {
				add(path)
			}
			return nil
		})
		return
	}
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		if e.Type().IsRegular() {
			add(filepath.Join(dir, e.Name()))
		}
	}
}

// addGlob adds regular files matching the glob pattern.
func addGlob(pat string, add func(string)) {
	matches, _ := filepath.Glob(pat)
	for _, m := range matches {
		if info, err := os.Stat(m); err == nil && info.Mode().IsRegular() {
			add(m)
		}
	}
}
