package main

import (
	"log"
	"net/http"

	"anonygram/internal/api"
	"anonygram/internal/config"
	"anonygram/internal/storage"
	"anonygram/internal/ws"
)

func main() {
	configs := config.Load()

	imageRepo := storage.NewLocalImageStore()

	fileRepo, err := storage.NewLocalFileStore(configs.UploadPath)
	if err != nil {
		log.Fatalf("Failed to initialize file store: %v", err)
	}

	hub := ws.NewHub(configs)
	go hub.Run()

	server := api.NewServer(imageRepo, fileRepo, configs, hub)

	handler := server.Routes()

	log.Printf("Starting server on :%s", configs.Port)
	log.Fatal(http.ListenAndServe(":"+configs.Port, handler))
}
