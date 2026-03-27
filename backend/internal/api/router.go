package api

import (
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

)

func (s *Server) Routes() http.Handler {
	r := chi.NewRouter()
	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   s.config.AllowedOrigins,
	}))

	// API routes
	r.Get("/images", s.ListImages)
	r.Post("/upload", s.UploadImage)

	// Static file serving
	fileServer := http.FileServer(http.Dir(s.config.UploadPath))
	r.Handle("/uploads/*", http.StripPrefix("/uploads/", fileServer))

	return r
}