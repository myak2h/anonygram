package storage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"anonygram/internal/models"
)

func TestNewLocalImageStore(t *testing.T) {
	store := NewLocalImageStore()
	assert.NotNil(t, store)
	assert.Empty(t, store.images)
}

func TestLocalImageStore_Add(t *testing.T) {
	store := NewLocalImageStore()

	img := models.Image{
		ID:        "test-id",
		Title:     "Test Image",
		Tags:      []string{"tag1", "tag2"},
		URL:       "/uploads/test-id.jpg",
		CreatedAt: time.Now(),
	}

	err := store.Add(img)

	assert.NoError(t, err)
	assert.Len(t, store.images, 1)
	assert.Equal(t, img, store.images[0])
}

func TestLocalImageStore_List(t *testing.T) {
	store := NewLocalImageStore()

	img1 := models.Image{
		ID:        "test-id-1",
		Title:     "Test Image 1",
		Tags:      []string{"tag1", "tag2"},
		URL:       "/uploads/test-id-1.jpg",
		CreatedAt: time.Now().Add(-time.Hour),
	}

	img2 := models.Image{
		ID:        "test-id-2",
		Title:     "Test Image 2",
		Tags:      []string{"tag2", "tag3"},
		URL:       "/uploads/test-id-2.jpg",
		CreatedAt: time.Now(),
	}

	err := store.Add(img1)
	assert.NoError(t, err)
	err = store.Add(img2)
	assert.NoError(t, err)

	t.Run("List all images", func(t *testing.T) {
		results := store.List(nil)
		assert.Len(t, results, 2)
		assert.Equal(t, img2, results[0])
		assert.Equal(t, img1, results[1])
	})

	t.Run("List images with tag2", func(t *testing.T) {
		results := store.List([]string{"tag2"})
		assert.Len(t, results, 2)
		assert.Equal(t, img2, results[0])
		assert.Equal(t, img1, results[1])
	})

	t.Run("List images with tag3", func(t *testing.T) {
		results := store.List([]string{"tag3"})
		assert.Len(t, results, 1)
		assert.Equal(t, img2, results[0])
	})

	t.Run("List images with non-existent tag", func(t *testing.T) {
		results := store.List([]string{"nonexistent"})
		assert.Len(t, results, 0)
	})
}

func TestLocalImageStore_List_Empty(t *testing.T) {
	store := NewLocalImageStore()
	results := store.List(nil)

	assert.NotNil(t, results, "should return empty slice, not nil")
	assert.Len(t, results, 0)
}
