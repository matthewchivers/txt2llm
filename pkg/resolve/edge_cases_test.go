package resolve

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFilesEdgeCases covers additional edge cases for better test coverage
func TestFilesEdgeCases(t *testing.T) {
	tmpDir := t.TempDir()

	// Test with a directory that exists but is empty
	emptyDir := filepath.Join(tmpDir, "empty")
	err := os.MkdirAll(emptyDir, 0750)
	assert.NoError(t, err)

	oldWd, err := os.Getwd()
	assert.NoError(t, err)
	assert.NoError(t, os.Chdir(tmpDir))
	defer func() { os.Chdir(oldWd) }()

	t.Run("empty directory non-recursive", func(t *testing.T) {
		files, err := Files([]string{"empty"}, false)
		assert.Error(t, err) // Should error because no files found
		assert.Nil(t, files)
		assert.Contains(t, err.Error(), "no files matched")
	})

	t.Run("empty directory recursive", func(t *testing.T) {
		files, err := Files([]string{"empty"}, true)
		assert.Error(t, err) // Should error because no files found
		assert.Nil(t, files)
		assert.Contains(t, err.Error(), "no files matched")
	})
}

// TestAddDirErrorHandling tests error conditions in addDir
func TestAddDirErrorHandling(t *testing.T) {
	t.Run("nonexistent directory", func(t *testing.T) {
		var collected []string
		add := func(path string) {
			if path != "" {
				collected = append(collected, path)
			}
		}

		// This should not panic and should gracefully handle the error
		addDir("/nonexistent/directory/path", false, add)

		// Should collect nothing
		assert.Empty(t, collected)
	})
}
