package models

import (
	"time"
)

// Feed represents an RSS channel
type Feed struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	ImageURL    string    `json:"image_url"`
	LastUpdated time.Time `json:"last_updated"`
	Items       []Item    `json:"items"`
}

// Item represents an individual news/article from an RSS channel
type Item struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Content     string    `json:"content"`
	Link        string    `json:"link"`
	Published   time.Time `json:"published"`
	Read        bool      `json:"read"`
}
