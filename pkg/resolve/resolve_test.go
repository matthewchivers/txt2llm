package resolve

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFiles(t *testing.T) {
	// Create temporary directory structure for testing
	tmpDir := t.TempDir()

	// Create test files and directories
	testFiles := map[string]string{
		"file1.txt":         "content1",
		"file2.go":          "package main",
		"subdir/file3.txt":  "content3",
		"subdir/file4.go":   "package sub",
		"subdir2/file5.md":  "# Header",
		"emptydir/.gitkeep": "",
	}

	for path, content := range testFiles {
		fullPath := filepath.Join(tmpDir, path)
		require.NoError(t, os.MkdirAll(filepath.Dir(fullPath), 0755))
		require.NoError(t, os.WriteFile(fullPath, []byte(content), 0644))
	}

	// Change to temp directory for relative path testing
	oldWd, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(tmpDir))
	defer func() { os.Chdir(oldWd) }()

	tests := []struct {
		name        string
		patterns    []string
		recursive   bool
		wantCount   int
		wantErr     bool
		contains    []string // files that should be in the result
		notContains []string // files that should not be in the result
	}{
		{
			name:      "single file",
			patterns:  []string{"file1.txt"},
			recursive: false,
			wantCount: 1,
			contains:  []string{"file1.txt"},
		},
		{
			name:      "multiple files",
			patterns:  []string{"file1.txt", "file2.go"},
			recursive: false,
			wantCount: 2,
			contains:  []string{"file1.txt", "file2.go"},
		},
		{
			name:        "directory non-recursive",
			patterns:    []string{"subdir"},
			recursive:   false,
			wantCount:   2,
			contains:    []string{"subdir/file3.txt", "subdir/file4.go"},
			notContains: []string{"file1.txt", "subdir2/file5.md"},
		},
		{
			name:      "directory recursive",
			patterns:  []string{"."},
			recursive: true,
			wantCount: 6, // all files including .gitkeep
			contains:  []string{"file1.txt", "file2.go", "subdir/file3.txt", "subdir/file4.go", "subdir2/file5.md", "emptydir/.gitkeep"},
		},
		{
			name:        "glob pattern",
			patterns:    []string{"*.txt"},
			recursive:   false,
			wantCount:   1,
			contains:    []string{"file1.txt"},
			notContains: []string{"subdir/file3.txt"}, // not in root
		},
		{
			name:      "glob pattern go files",
			patterns:  []string{"*.go"},
			recursive: false,
			wantCount: 1,
			contains:  []string{"file2.go"},
		},
		{
			name:      "multiple patterns",
			patterns:  []string{"*.txt", "*.go"},
			recursive: false,
			wantCount: 2,
			contains:  []string{"file1.txt", "file2.go"},
		},
		{
			name:      "nonexistent file",
			patterns:  []string{"nonexistent.txt"},
			recursive: false,
			wantErr:   true,
		},
		{
			name:      "empty patterns",
			patterns:  []string{},
			recursive: false,
			wantErr:   true,
		},
		{
			name:      "empty pattern strings",
			patterns:  []string{"", "file1.txt", ""},
			recursive: false,
			wantCount: 1,
			contains:  []string{"file1.txt"},
		},
		{
			name:      "duplicate files",
			patterns:  []string{"file1.txt", "file1.txt"},
			recursive: false,
			wantCount: 1, // should be deduplicated
			contains:  []string{"file1.txt"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files, err := Files(tt.patterns, tt.recursive)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, files)
				return
			}

			require.NoError(t, err)
			assert.Len(t, files, tt.wantCount)

			// Convert to absolute paths for comparison
			var absFiles []string
			for _, f := range files {
				abs, _ := filepath.Abs(f)
				absFiles = append(absFiles, abs)
			}

			// Check that expected files are present
			for _, expectedFile := range tt.contains {
				expectedAbs, _ := filepath.Abs(expectedFile)
				assert.Contains(t, absFiles, expectedAbs, "expected file %s not found", expectedFile)
			}

			// Check that unexpected files are not present
			for _, unexpectedFile := range tt.notContains {
				unexpectedAbs, _ := filepath.Abs(unexpectedFile)
				assert.NotContains(t, absFiles, unexpectedAbs, "unexpected file %s found", unexpectedFile)
			}

			// Ensure all returned paths are absolute
			for _, file := range files {
				assert.True(t, filepath.IsAbs(file), "path should be absolute: %s", file)
			}

			// Ensure no duplicates
			seen := make(map[string]bool)
			for _, file := range files {
				assert.False(t, seen[file], "duplicate file found: %s", file)
				seen[file] = true
			}
		})
	}
}

func TestAddDir(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test structure
	testFiles := map[string]string{
		"root1.txt":          "content",
		"root2.go":           "package main",
		"sub1/file1.txt":     "content1",
		"sub1/file2.go":      "package sub1",
		"sub1/sub2/file3.md": "markdown",
	}

	for path, content := range testFiles {
		fullPath := filepath.Join(tmpDir, path)
		require.NoError(t, os.MkdirAll(filepath.Dir(fullPath), 0755))
		require.NoError(t, os.WriteFile(fullPath, []byte(content), 0644))
	}

	tests := []struct {
		name      string
		dir       string
		recursive bool
		expected  []string
	}{
		{
			name:      "non-recursive",
			dir:       tmpDir,
			recursive: false,
			expected:  []string{"root1.txt", "root2.go"}, // only root level files
		},
		{
			name:      "recursive",
			dir:       tmpDir,
			recursive: true,
			expected:  []string{"root1.txt", "root2.go", "sub1/file1.txt", "sub1/file2.go", "sub1/sub2/file3.md"},
		},
		{
			name:      "subdirectory non-recursive",
			dir:       filepath.Join(tmpDir, "sub1"),
			recursive: false,
			expected:  []string{"file1.txt", "file2.go"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var collected []string
			add := func(path string) {
				if path != "" {
					// Make path relative to the directory being scanned for easier testing
					rel, err := filepath.Rel(tt.dir, path)
					if err == nil && !filepath.IsAbs(rel) {
						collected = append(collected, rel)
					}
				}
			}

			addDir(tt.dir, tt.recursive, add)

			assert.ElementsMatch(t, tt.expected, collected)
		})
	}
}

func TestAddGlob(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files
	testFiles := map[string]string{
		"file1.txt": "content1",
		"file2.go":  "package main",
		"file3.md":  "# Header",
		"file4.txt": "content4",
	}

	for path, content := range testFiles {
		fullPath := filepath.Join(tmpDir, path)
		require.NoError(t, os.WriteFile(fullPath, []byte(content), 0644))
	}

	oldWd, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(tmpDir))
	defer func() { os.Chdir(oldWd) }()

	tests := []struct {
		name     string
		pattern  string
		expected []string
	}{
		{
			name:     "txt files",
			pattern:  "*.txt",
			expected: []string{"file1.txt", "file4.txt"},
		},
		{
			name:     "go files",
			pattern:  "*.go",
			expected: []string{"file2.go"},
		},
		{
			name:     "md files",
			pattern:  "*.md",
			expected: []string{"file3.md"},
		},
		{
			name:     "all files",
			pattern:  "*",
			expected: []string{"file1.txt", "file2.go", "file3.md", "file4.txt"},
		},
		{
			name:     "no matches",
			pattern:  "*.xyz",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var collected []string
			add := func(path string) {
				if path != "" {
					collected = append(collected, path)
				}
			}

			addGlob(tt.pattern, add)

			assert.ElementsMatch(t, tt.expected, collected)
		})
	}
}
