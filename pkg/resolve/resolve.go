// Package resolve handles file pattern resolution for txt2llm.
package resolve

import (
	"fmt"
	"os"
	"path/filepath"
)

// Files resolves patterns (files, directories, globs) to a deduplicated slice of
// absolute file paths. Returns an error if no files match.
func Files(patterns []string, recursive bool) ([]string, error) {
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
	if len(out) == 0 {
		return nil, fmt.Errorf("no files matched any of the patterns: %v", patterns)
	}
	return out, nil
}

// addDir adds regular files from the specified directory.
// If recursive is true, it walks the directory tree.
func addDir(dir string, recursive bool, add func(string)) {
	if recursive {
		_ = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
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
