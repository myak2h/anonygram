package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"

	"anonygram/internal/config"
	"anonygram/internal/models"
	"anonygram/internal/storage"
	"anonygram/internal/utils"
	"anonygram/internal/ws"
)

type Server struct {
	imageRepo storage.ImageRepository
	fileRepo  storage.FileRepository
	config    *config.Config
	hub       ws.Broadcaster
}

func NewServer(imageRepo storage.ImageRepository, fileRepo storage.FileRepository, config *config.Config, hub ws.Broadcaster) *Server {
	return &Server{
		imageRepo: imageRepo,
		fileRepo:  fileRepo,
		config:    config,
		hub:       hub,
	}
}

// Handlers

func (s *Server) ListImages(w http.ResponseWriter, r *http.Request) {
	tags := r.URL.Query()["tag"]
	images := s.imageRepo.List(tags)
	respondWithJSON(w, http.StatusOK, images)
}

func (s *Server) UploadImage(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(s.config.MaxUploadSize); err != nil {
		if err.Error() == "http: request body too large" {
			respondWithError(w, http.StatusRequestEntityTooLarge, ErrFileTooLarge.Error())
			return
		}
		respondWithError(w, http.StatusBadRequest, ErrInvalidFormData.Error())
		return
	}

	title := r.FormValue("title")
	if title == "" {
		respondWithError(w, http.StatusBadRequest, ErrTitleRequired.Error())
		return
	}

	tags := utils.SplitAndTrim(r.FormValue("tags"), ",")

	file, _, err := r.FormFile("image")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, ErrImageFileRequired.Error())
		return
	}
	defer func() { _ = file.Close() }()

	filename, err := s.fileRepo.Save(file)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, ErrFileSaveFailed.Error())
		return
	}

	img := models.Image{
		ID:        uuid.New().String(),
		Title:     title,
		Tags:      tags,
		Filename:  filename,
		CreatedAt: time.Now(),
	}

	if err := s.imageRepo.Add(img); err != nil {
		respondWithError(w, http.StatusInternalServerError, ErrMetadataSaveFailed.Error())
		return
	}

	if !s.hub.Broadcast(img) {
		log.Printf("Failed to broadcast new image: %s", img.ID)
	}

	respondWithJSON(w, http.StatusCreated, img)
}

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, status int, message string) {
	respondWithJSON(w, status, map[string]string{"error": message})
}
