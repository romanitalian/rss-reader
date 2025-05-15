package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/romanitalian/rss-reader/internal/models"
)

// Storage defines the interface for RSS feed storage
type Storage interface {
	SaveFeed(feed *models.Feed) error
	GetFeedByID(id string) (*models.Feed, error)
	GetAllFeeds() ([]*models.Feed, error)
	DeleteFeed(id string) error
	UpdateFeed(feed *models.Feed) error
	MarkItemAsRead(feedID, itemID string) error
}

// FileStorage represents a file system based storage implementation
type FileStorage struct {
	dataDir string
}

// NewFileStorage creates a new instance of FileStorage
func NewFileStorage(dataDir string) (*FileStorage, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}
	return &FileStorage{dataDir: dataDir}, nil
}

// SaveFeed saves an RSS feed to storage
func (s *FileStorage) SaveFeed(feed *models.Feed) error {
	if feed.ID == "" {
		return errors.New("feed ID cannot be empty")
	}

	feed.LastUpdated = time.Now()
	data, err := json.Marshal(feed)
	if err != nil {
		return err
	}

	filePath := filepath.Join(s.dataDir, feed.ID+".json")
	return os.WriteFile(filePath, data, 0644)
}

// GetFeedByID returns an RSS feed by its ID
func (s *FileStorage) GetFeedByID(id string) (*models.Feed, error) {
	filePath := filepath.Join(s.dataDir, id+".json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	feed := &models.Feed{}
	if err := json.Unmarshal(data, feed); err != nil {
		return nil, err
	}

	return feed, nil
}

// GetAllFeeds returns all saved RSS feeds
func (s *FileStorage) GetAllFeeds() ([]*models.Feed, error) {
	files, err := os.ReadDir(s.dataDir)
	if err != nil {
		return nil, err
	}

	feeds := make([]*models.Feed, 0, len(files))
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}

		data, err := os.ReadFile(filepath.Join(s.dataDir, file.Name()))
		if err != nil {
			continue
		}

		feed := &models.Feed{}
		if err := json.Unmarshal(data, feed); err != nil {
			continue
		}

		feeds = append(feeds, feed)
	}

	return feeds, nil
}

// DeleteFeed removes an RSS feed from storage
func (s *FileStorage) DeleteFeed(id string) error {
	filePath := filepath.Join(s.dataDir, id+".json")
	return os.Remove(filePath)
}

// UpdateFeed updates an RSS feed in storage
func (s *FileStorage) UpdateFeed(feed *models.Feed) error {
	return s.SaveFeed(feed)
}

// MarkItemAsRead marks an article as read
func (s *FileStorage) MarkItemAsRead(feedID, itemID string) error {
	feed, err := s.GetFeedByID(feedID)
	if err != nil {
		return err
	}

	found := false
	for i := range feed.Items {
		if feed.Items[i].ID == itemID {
			feed.Items[i].Read = true
			found = true
			break
		}
	}

	if !found {
		return errors.New("item not found")
	}

	return s.SaveFeed(feed)
}
