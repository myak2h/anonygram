package storage

import (
	"anonygram/internal/models"
	"io"
)

type ImageRepository interface {
	Add(img models.Image) error
	List(tags []string) []models.Image
}

type FileRepository interface {
	Save(src io.Reader) (string, error)
}