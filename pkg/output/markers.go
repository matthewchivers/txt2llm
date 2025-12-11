// Package output handles file output formatting with markers for txt2llm.
package output

import (
	"fmt"
	"os"
	"path/filepath"
)

// Header prints a concise explanation of markers.
func Header(markerPrefix, markerSuffix string) {
	fmt.Printf("Each section below represents text output from one file.\n")
	fmt.Printf("Delimiters: %sSTART:{filename}%s ... %sEND:{filename}%s\n\n", markerPrefix, markerSuffix, markerPrefix, markerSuffix)
}

// Paths returns either absolute or relative paths depending on flag.
func Paths(files []string, relative bool) []string {
	if !relative {
		return files
	}
	cwd, _ := os.Getwd()
	rel := make([]string, 0, len(files))
	for _, p := range files {
		if r, err := filepath.Rel(cwd, p); err == nil {
			rel = append(rel, r)
		} else {
			rel = append(rel, p)
		}
	}
	return rel
}

// Markers emits all files with start/end markers.
func Markers(files []string, outPaths []string, markerPrefix, markerSuffix string) {
	for i, src := range files {
		emit(src, outPaths[i], markerPrefix, markerSuffix)
	}
}

func emit(srcPath, outPath, markerPrefix, markerSuffix string) {
	data, err := os.ReadFile(srcPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", srcPath, err)
		return
	}
	fmt.Printf("%sSTART:%s%s\n", markerPrefix, outPath, markerSuffix)
	_, _ = os.Stdout.Write(data) // Ignore write errors to stdout
	newlineIfNeeded(data)
	fmt.Printf("%sEND:%s%s\n\n", markerPrefix, outPath, markerSuffix)
}

func newlineIfNeeded(data []byte) {
	if len(data) == 0 || data[len(data)-1] == '\n' {
		return
	}
	fmt.Println()
}
