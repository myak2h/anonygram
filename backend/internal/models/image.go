package models

import "time"

type Image struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Tags      []string  `json:"tags"`
	Filename  string    `json:"filename"`
	CreatedAt time.Time `json:"createdAt"`
}
