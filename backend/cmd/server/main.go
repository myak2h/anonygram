package main

import (
	"log"
	"net/http"

	"anonygram/internal/api"
	"anonygram/internal/config"
	"anonygram/internal/storage"
)

func main() {
	configs := config.Load()

	imageRepo := storage.NewLocalImageStore()

	fileRepo, err := storage.NewLocalFileStore(configs.UploadPath)
	if err != nil {
		log.Fatalf("Failed to initialize file store: %v", err)
	}

	server := api.NewServer(imageRepo, fileRepo, configs)

	handler := server.Routes()

	log.Printf("Starting server on :%s", configs.Port)
	log.Fatal(http.ListenAndServe(":"+configs.Port, handler))
}