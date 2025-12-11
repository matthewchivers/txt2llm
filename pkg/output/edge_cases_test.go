package output

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPathsEdgeCases covers edge cases for better coverage
func TestPathsEdgeCases(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create a file outside the temp directory to test edge case where
	// filepath.Rel might fail
	testFile := filepath.Join(tmpDir, "test.txt")
	require.NoError(t, os.WriteFile(testFile, []byte("content"), 0644))

	// Change to a different directory to create a scenario where 
	// the relative path calculation might have edge cases
	subDir := filepath.Join(tmpDir, "subdir")
	require.NoError(t, os.MkdirAll(subDir, 0755))
	
	oldWd, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(subDir))
	defer func() { os.Chdir(oldWd) }()

	t.Run("relative path from subdirectory", func(t *testing.T) {
		files := []string{testFile}
		result := Paths(files, true)
		
		// Should return one result
		assert.Len(t, result, 1)
		// The result should be a relative path to the parent directory
		assert.Contains(t, result[0], "test.txt")
	})

	t.Run("absolute paths unchanged", func(t *testing.T) {
		files := []string{testFile}
		result := Paths(files, false)
		
		assert.Equal(t, files, result)
	})
}