package api

import "errors"

var (
	ErrTitleRequired      = errors.New("title is required")
	ErrImageFileRequired  = errors.New("image file is required")
	ErrInvalidFormData    = errors.New("invalid form data")
	ErrFileTooLarge       = errors.New("file is too large")
	ErrFileSaveFailed     = errors.New("failed to save image file")
	ErrMetadataSaveFailed = errors.New("failed to save image metadata")
)
