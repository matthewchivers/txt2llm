package output

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeader(t *testing.T) {
	tests := []struct {
		name         string
		markerPrefix string
		markerSuffix string
		expected     string
	}{
		{
			name:         "default markers",
			markerPrefix: "<<<",
			markerSuffix: ">>>",
			expected:     "Each section below represents text output from one file.\nDelimiters: <<<START:{filename}>>> ... <<<END:{filename}>>>\n\n",
		},
		{
			name:         "custom markers",
			markerPrefix: "[[[",
			markerSuffix: "]]]",
			expected:     "Each section below represents text output from one file.\nDelimiters: [[[START:{filename}]]] ... [[[END:{filename}]]]\n\n",
		},
		{
			name:         "empty markers",
			markerPrefix: "",
			markerSuffix: "",
			expected:     "Each section below represents text output from one file.\nDelimiters: START:{filename} ... END:{filename}\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			Header(tt.markerPrefix, tt.markerSuffix)

			w.Close()
			os.Stdout = old

			buf := make([]byte, 1024)
			n, _ := r.Read(buf)
			output := string(buf[:n])

			assert.Equal(t, tt.expected, output)
		})
	}
}

func TestPaths(t *testing.T) {
	// Create temporary directory for testing
	tmpDir := t.TempDir()

	// Create some test files
	file1 := filepath.Join(tmpDir, "file1.txt")
	file2 := filepath.Join(tmpDir, "subdir", "file2.txt")
	require.NoError(t, os.MkdirAll(filepath.Dir(file2), 0755))
	require.NoError(t, os.WriteFile(file1, []byte("content"), 0644))
	require.NoError(t, os.WriteFile(file2, []byte("content"), 0644))

	// Change to temp directory
	oldWd, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(tmpDir))
	defer func() { os.Chdir(oldWd) }()

	tests := []struct {
		name     string
		files    []string
		relative bool
		expected []string
	}{
		{
			name:     "absolute paths when relative=false",
			files:    []string{file1, file2},
			relative: false,
			expected: []string{file1, file2},
		},
		{
			name:     "relative paths when relative=true",
			files:    []string{file1, file2},
			relative: true,
			expected: []string{"file1.txt", "subdir/file2.txt"},
		},
		{
			name:     "empty slice",
			files:    []string{},
			relative: false,
			expected: []string{},
		},
		{
			name:     "empty slice relative",
			files:    []string{},
			relative: true,
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Paths(tt.files, tt.relative)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPathsWithInvalidPaths(t *testing.T) {
	// Test with paths that can't be made relative
	invalidPath := "/some/absolute/path/that/does/not/exist"
	files := []string{invalidPath}

	result := Paths(files, true)

	// Should have one result, either original or relative version
	assert.Len(t, result, 1)
	// The result should be some version of the path (either original or relative)
	assert.Contains(t, result[0], "some/absolute/path/that/does/not/exist")
}

func TestMarkers(t *testing.T) {
	// Create temporary files for testing
	tmpDir := t.TempDir()

	file1 := filepath.Join(tmpDir, "file1.txt")
	file2 := filepath.Join(tmpDir, "file2.go")

	content1 := "Hello, World!"
	content2 := "package main\n\nfunc main() {\n\tprintln(\"test\")\n}"

	require.NoError(t, os.WriteFile(file1, []byte(content1), 0644))
	require.NoError(t, os.WriteFile(file2, []byte(content2), 0644))

	tests := []struct {
		name             string
		files            []string
		outPaths         []string
		markerPrefix     string
		markerSuffix     string
		expectedContains []string
	}{
		{
			name:         "single file",
			files:        []string{file1},
			outPaths:     []string{"file1.txt"},
			markerPrefix: "<<<",
			markerSuffix: ">>>",
			expectedContains: []string{
				"<<<START:file1.txt>>>",
				content1,
				"<<<END:file1.txt>>>",
			},
		},
		{
			name:         "multiple files",
			files:        []string{file1, file2},
			outPaths:     []string{"file1.txt", "file2.go"},
			markerPrefix: "[[[",
			markerSuffix: "]]]",
			expectedContains: []string{
				"[[[START:file1.txt]]]",
				content1,
				"[[[END:file1.txt]]]",
				"[[[START:file2.go]]]",
				content2,
				"[[[END:file2.go]]]",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			Markers(tt.files, tt.outPaths, tt.markerPrefix, tt.markerSuffix)

			w.Close()
			os.Stdout = old

			buf := make([]byte, 4096)
			n, _ := r.Read(buf)
			output := string(buf[:n])

			// Check that all expected content is present
			for _, expected := range tt.expectedContains {
				assert.Contains(t, output, expected, "Output should contain: %s", expected)
			}
		})
	}
}

func TestEmit(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name             string
		content          string
		outPath          string
		markerPrefix     string
		markerSuffix     string
		expectedContains []string
	}{
		{
			name:         "file with newline",
			content:      "Hello\nWorld\n",
			outPath:      "test.txt",
			markerPrefix: "<<<",
			markerSuffix: ">>>",
			expectedContains: []string{
				"<<<START:test.txt>>>",
				"Hello\nWorld\n",
				"<<<END:test.txt>>>",
			},
		},
		{
			name:         "file without newline",
			content:      "Hello World",
			outPath:      "test.txt",
			markerPrefix: "<<<",
			markerSuffix: ">>>",
			expectedContains: []string{
				"<<<START:test.txt>>>",
				"Hello World",
				"<<<END:test.txt>>>",
			},
		},
		{
			name:         "empty file",
			content:      "",
			outPath:      "empty.txt",
			markerPrefix: "===",
			markerSuffix: "===",
			expectedContains: []string{
				"===START:empty.txt===",
				"===END:empty.txt===",
			},
		},
		{
			name:         "binary-like content",
			content:      "binary\x00\x01\x02content",
			outPath:      "binary.dat",
			markerPrefix: "---",
			markerSuffix: "---",
			expectedContains: []string{
				"---START:binary.dat---",
				"binary\x00\x01\x02content",
				"---END:binary.dat---",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test file
			testFile := filepath.Join(tmpDir, tt.name+".txt")
			require.NoError(t, os.WriteFile(testFile, []byte(tt.content), 0644))

			// Capture stdout
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			emit(testFile, tt.outPath, tt.markerPrefix, tt.markerSuffix)

			w.Close()
			os.Stdout = old

			buf := make([]byte, 4096)
			n, _ := r.Read(buf)
			output := string(buf[:n])

			// Check expected content
			for _, expected := range tt.expectedContains {
				assert.Contains(t, output, expected)
			}
		})
	}
}

func TestEmitWithNonexistentFile(t *testing.T) {
	// Capture stderr
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	emit("nonexistent.txt", "nonexistent.txt", "<<<", ">>>")

	w.Close()
	os.Stderr = old

	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	output := string(buf[:n])

	assert.Contains(t, output, "Error reading nonexistent.txt")
}

func TestNewlineIfNeeded(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected string
	}{
		{
			name:     "ends with newline",
			data:     []byte("hello\n"),
			expected: "",
		},
		{
			name:     "no newline",
			data:     []byte("hello"),
			expected: "\n",
		},
		{
			name:     "empty data",
			data:     []byte(""),
			expected: "",
		},
		{
			name:     "multiple newlines",
			data:     []byte("hello\n\n"),
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			newlineIfNeeded(tt.data)

			w.Close()
			os.Stdout = old

			var buf bytes.Buffer
			buf.ReadFrom(r)
			output := buf.String()

			assert.Equal(t, tt.expected, output)
		})
	}
}
