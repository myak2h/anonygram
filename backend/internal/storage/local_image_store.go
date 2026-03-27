package storage

import (
	"sort"
	"strings"
	"sync"

	"anonygram/internal/models"
)

type LocalImageStore struct {
	mu     sync.RWMutex
	images []models.Image
}

func NewLocalImageStore() *LocalImageStore {
	return &LocalImageStore{
		images: make([]models.Image, 0),
	}
}

func (s *LocalImageStore) Add(img models.Image) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.images = append(s.images, img)
	return nil
}

func (s *LocalImageStore) List(tags []string) []models.Image {
	s.mu.RLock()
	if len(s.images) == 0 {
		s.mu.RUnlock()
		return make([]models.Image, 0)
	}
	imagesCopy := make([]models.Image, len(s.images))
	copy(imagesCopy, s.images)
	s.mu.RUnlock()

	if len(tags) == 0 {
		sortImages(imagesCopy)
		return imagesCopy
	}

	tagSet := make(map[string]struct{}, len(tags))
	for _, tag := range tags {
		tagSet[normalizeTag(tag)] = struct{}{}
	}

	filteredImages := make([]models.Image, 0)

	for _, img := range imagesCopy {
		for _, tag := range img.Tags {
			if _, ok := tagSet[normalizeTag(tag)]; ok {
				filteredImages = append(filteredImages, img)
				break
			}
		}
	}

	sortImages(filteredImages)
	return filteredImages
}

func normalizeTag(tag string) string {
	return strings.ToLower(strings.TrimSpace(tag))
}

func sortImages(images []models.Image) {
	sort.Slice(images, func(i, j int) bool {
		return images[i].CreatedAt.After(images[j].CreatedAt)
	})
}
