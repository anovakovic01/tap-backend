package news

import (
	"time"
)

// Repository contains API for storing and reading News data.
type Repository interface {
	// Creates new News instance in data store.
	Create(News) (int64, error)

	// One returns News instance by ID if it exists.
	One(int64) (News, error)

	// All returns list of all News stored in data store.
	All() []News
}

// Collector collects list of news.
type Collector interface {
	// Collects list of news from external sources.
	Collect(string) ([]News, error)
}

// News contains new related data.
type News struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Link        string    `json:"link"`
	Description string    `json:"description"`
	ImageTitle  string    `json:"imageTitle"`
	Image       string    `json:"image"`
	PubDate     time.Time `json:"pubDate"`
}
