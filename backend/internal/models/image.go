package models

import "time"

type Image struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Tags      []string  `json:"tags"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"createdAt"`
}