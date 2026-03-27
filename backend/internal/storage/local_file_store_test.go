package storage

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLocalFileStore(t *testing.T) {
	tempDir := t.TempDir()
	basePath := filepath.Join(tempDir, "uploads")
	store, err := NewLocalFileStore(basePath)

	assert.NoError(t, err)
	assert.NotNil(t, store)

	info, err := os.Stat(basePath)
	assert.NoError(t, err)
	assert.True(t, info.IsDir())
}

func TestLocalFileStore_Save(t *testing.T) {
	tempDir := t.TempDir()
	fileStore, err := NewLocalFileStore(tempDir)
	assert.NoError(t, err)

	t.Run("Save PNG file", func(t *testing.T) {

		// PNG file header bytes
		content := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52}
		src := bytes.NewReader(content)

		path, err := fileStore.Save(src)
		assert.NoError(t, err)
		assert.Equal(t, ".png", filepath.Ext(path))
		assert.FileExists(t, path)

		data, err := os.ReadFile(path)
		assert.NoError(t, err)
		assert.Equal(t, content, data)
	})

	t.Run("Save JPEG file", func(t *testing.T) {
		// JPEG file header bytes
		content := []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46, 0x00, 0x01}
		src := bytes.NewReader(content)

		path, err := fileStore.Save(src)
		assert.NoError(t, err)
		assert.Equal(t, ".jpg", filepath.Ext(path))
		assert.FileExists(t, path)

		data, err := os.ReadFile(path)
		assert.NoError(t, err)
		assert.Equal(t, content, data)
	})

	t.Run("Save GIF file", func(t *testing.T) {
		// GIF file header bytes
		content := []byte{0x47, 0x49, 0x46, 0x38, 0x39, 0x61, 0x01, 0x00, 0x01, 0x00}
		src := bytes.NewReader(content)

		path, err := fileStore.Save(src)
		assert.NoError(t, err)
		assert.Equal(t, ".gif", filepath.Ext(path))
		assert.FileExists(t, path)

		data, err := os.ReadFile(path)
		assert.NoError(t, err)
		assert.Equal(t, content, data)
	})

	t.Run("Save unsupported file type", func(t *testing.T) {
		content := []byte{0x00, 0x01, 0x02}
		src := bytes.NewReader(content)

		path, err := fileStore.Save(src)
		assert.Error(t, err)
		assert.Empty(t, path)
		assert.Contains(t, err.Error(), "unsupported file type")
	})
}
