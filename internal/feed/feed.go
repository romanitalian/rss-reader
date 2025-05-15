package feed

import (
	"crypto/md5"
	"fmt"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/romanitalian/rss-reader/internal/models"
)

// FeedService provides functions for working with RSS feeds
type FeedService struct {
	parser *gofeed.Parser
}

// NewFeedService creates a new instance of FeedService
func NewFeedService() *FeedService {
	return &FeedService{
		parser: gofeed.NewParser(),
	}
}

// FetchFeed downloads an RSS feed from the specified URL
func (s *FeedService) FetchFeed(url string) (*models.Feed, error) {
	gfeed, err := s.parser.ParseURL(url)
	if err != nil {
		return nil, err
	}

	feedID := generateID(url)
	feed := &models.Feed{
		ID:          feedID,
		Title:       gfeed.Title,
		Description: gfeed.Description,
		URL:         url,
		LastUpdated: time.Now(),
		Items:       make([]models.Item, 0, len(gfeed.Items)),
	}

	if gfeed.Image != nil {
		feed.ImageURL = gfeed.Image.URL
	}

	for _, gitem := range gfeed.Items {
		item := models.Item{
			ID:          generateID(gitem.GUID + gitem.Link),
			Title:       gitem.Title,
			Description: gitem.Description,
			Link:        gitem.Link,
			Read:        false,
		}

		if gitem.Content != "" {
			item.Content = gitem.Content
		} else {
			item.Content = gitem.Description
		}

		if gitem.PublishedParsed != nil {
			item.Published = *gitem.PublishedParsed
		} else if gitem.UpdatedParsed != nil {
			item.Published = *gitem.UpdatedParsed
		} else {
			item.Published = time.Now()
		}

		feed.Items = append(feed.Items, item)
	}

	return feed, nil
}

// UpdateFeed updates the feed while preserving read articles
func (s *FeedService) UpdateFeed(oldFeed *models.Feed) (*models.Feed, error) {
	newFeed, err := s.FetchFeed(oldFeed.URL)
	if err != nil {
		return nil, err
	}

	// Save the state of read articles
	readItems := make(map[string]bool)
	for _, item := range oldFeed.Items {
		if item.Read {
			readItems[item.ID] = true
		}
	}

	// Check if there are already articles with the same ID and mark them as read if needed
	for i, item := range newFeed.Items {
		if _, exists := readItems[item.ID]; exists {
			newFeed.Items[i].Read = true
		}
	}

	return newFeed, nil
}

// generateID generates a unique ID based on the input string
func generateID(input string) string {
	hash := md5.Sum([]byte(input))
	return fmt.Sprintf("%x", hash)
}
