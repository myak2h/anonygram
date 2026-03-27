package api

import "errors"


var (
	ErrTitleRequired      = errors.New("Title is required")
	ErrImageFileRequired  = errors.New("Image file is required")
	ErrInvalidFormData    = errors.New("Invalid form data")
	ErrFileTooLarge       = errors.New("File is too large")
	ErrFileSaveFailed     = errors.New("Failed to save image file")
	ErrMetadataSaveFailed = errors.New("Failed to save image metadata")
)


