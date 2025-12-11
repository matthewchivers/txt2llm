package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMainUnit verifies the main function correctly processes command-line arguments and outputs files with proper markers.
func TestMainUnit(t *testing.T) {
	// Create a test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	require.NoError(t, os.WriteFile(testFile, []byte("test content"), 0644))

	// Change to the temp directory
	oldWd, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(tmpDir))
	defer func() { os.Chdir(oldWd) }()

	// Test successful execution
	t.Run("successful execution", func(t *testing.T) {
		// Save original args, stdout, and command line
		oldArgs := os.Args
		old := os.Stdout
		oldCmdLine := pflag.CommandLine
		defer func() {
			os.Args = oldArgs
			os.Stdout = old
			pflag.CommandLine = oldCmdLine
		}()

		// Reset pflag state
		pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

		// Mock command line args
		os.Args = []string{"txt2llm", "test.txt"}

		// Capture stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		// Call main
		main()

		// Close and restore
		w.Close()
		os.Stdout = old

		// Read output
		buf := make([]byte, 1024)
		n, _ := r.Read(buf)
		output := string(buf[:n])

		// Verify output contains expected markers and content
		assert.Contains(t, output, "<<<START:")
		assert.Contains(t, output, "test.txt>>>")
		assert.Contains(t, output, "test content")
		assert.Contains(t, output, "<<<END:")
	})

	t.Run("with relative paths", func(t *testing.T) {
		// Save original args, stdout, and command line
		oldArgs := os.Args
		old := os.Stdout
		oldCmdLine := pflag.CommandLine
		defer func() {
			os.Args = oldArgs
			os.Stdout = old
			pflag.CommandLine = oldCmdLine
		}()

		// Reset pflag state
		pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

		// Mock command line args with relative flag
		os.Args = []string{"txt2llm", "--relative", "test.txt"}

		// Capture stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		// Call main
		main()

		// Close and restore
		w.Close()
		os.Stdout = old

		// Read output
		buf := make([]byte, 1024)
		n, _ := r.Read(buf)
		output := string(buf[:n])

		// Verify output contains expected markers and content with relative path
		assert.Contains(t, output, "<<<START:test.txt>>>")
		assert.Contains(t, output, "test content")
		assert.Contains(t, output, "<<<END:test.txt>>>")
	})

	t.Run("with custom markers", func(t *testing.T) {
		// Save original args, stdout, and command line
		oldArgs := os.Args
		old := os.Stdout
		oldCmdLine := pflag.CommandLine
		defer func() {
			os.Args = oldArgs
			os.Stdout = old
			pflag.CommandLine = oldCmdLine
		}()

		// Reset pflag state
		pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

		// Mock command line args with custom markers
		os.Args = []string{"txt2llm", "--marker-prefix", "[[[", "--marker-suffix", "]]]", "test.txt"}

		// Capture stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		// Call main
		main()

		// Close and restore
		w.Close()
		os.Stdout = old

		// Read output
		buf := make([]byte, 1024)
		n, _ := r.Read(buf)
		output := string(buf[:n])

		// Verify output contains expected custom markers
		assert.Contains(t, output, "[[[START:")
		assert.Contains(t, output, "test.txt]]]")
		assert.Contains(t, output, "test content")
		assert.Contains(t, output, "[[[END:")
	})

	t.Run("with multiple files", func(t *testing.T) {
		// Create another test file
		testFile2 := filepath.Join(tmpDir, "test2.go")
		require.NoError(t, os.WriteFile(testFile2, []byte("package main"), 0644))

		// Save original args, stdout, and command line
		oldArgs := os.Args
		old := os.Stdout
		oldCmdLine := pflag.CommandLine
		defer func() {
			os.Args = oldArgs
			os.Stdout = old
			pflag.CommandLine = oldCmdLine
		}()

		// Reset pflag state
		pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

		// Mock command line args with multiple files
		os.Args = []string{"txt2llm", "--relative", "test.txt", "test2.go"}

		// Capture stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		// Call main
		main()

		// Close and restore
		w.Close()
		os.Stdout = old

		// Read output
		buf := make([]byte, 2048)
		n, _ := r.Read(buf)
		output := string(buf[:n])

		// Verify both files are processed
		assert.Contains(t, output, "<<<START:test.txt>>>")
		assert.Contains(t, output, "test content")
		assert.Contains(t, output, "<<<END:test.txt>>>")
		assert.Contains(t, output, "<<<START:test2.go>>>")
		assert.Contains(t, output, "package main")
		assert.Contains(t, output, "<<<END:test2.go>>>")
	})
}
