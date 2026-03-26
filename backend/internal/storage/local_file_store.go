package storage

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)


type LocalFileStore struct {
	basePath string
}


func NewLocalFileStore(basePath string) (*LocalFileStore, error) {
	if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
		return nil, err
	}

	return &LocalFileStore{
		basePath: basePath,
	}, nil
}

func (s *LocalFileStore) Save(src io.Reader) (string, error) {
	head, ext, err := detectFileType(src)
	if err != nil {
		return "", err
	}

	filename := uuid.NewString() + ext
	destPath := filepath.Join(s.basePath, filename)
	
	destFile, err := os.Create(destPath)
	if err != nil {
		return "", err
	}

	defer destFile.Close()

	fullReader := io.MultiReader(bytes.NewReader(head), src)
	
	if _, err := io.Copy(destFile, fullReader); err != nil {
		destFile.Close()
		os.Remove(destPath)
		return "", err
	}

	return destPath, nil
}

func detectFileType(src io.Reader) ([]byte, string, error) {
	head := make([]byte, 512)

	n, err := io.ReadFull(src, head)
	if err != nil && err != io.ErrUnexpectedEOF {
		return nil, "", err
	}

	head = head[:n]

	contentType := http.DetectContentType(head)

	switch contentType {
	case "image/jpeg":
		return head, ".jpg", nil
	case "image/png":
		return head, ".png", nil
	case "image/gif":
		return head, ".gif", nil
	default:
		return nil, "", fmt.Errorf("unsupported file type: %s", contentType)
	}
}

var _ FileRepository = (*LocalFileStore)(nil)
