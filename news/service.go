package news

import (
	"errors"
	"fmt"
	"time"
)

// ErrNotFound indicates that news doesn't exist
var ErrNotFound = errors.New("news not found")

var _ Service = (*newsService)(nil)

// Service contains business logic API related to News.
type Service interface {
	// One returns News instance by ID if it exists.
	One(int64) (News, error)

	// All returns list of all News.
	All() []News

	// Collects list of news.
	Collect(string)
}

type newsService struct {
	repo      Repository
	collector Collector
}

// NewService instantiates new news service.
func NewService(repo Repository, collector Collector) Service {
	return newsService{repo, collector}
}

func (svc newsService) One(id int64) (News, error) {
	item, err := svc.repo.One(id)
	if err != nil {
		return News{}, ErrNotFound
	}

	return item, nil
}

func (svc newsService) All() []News {
	return svc.repo.All()
}

func (svc newsService) Collect(src string) {
	for {
		items, err := svc.collector.Collect(src)
		if err != nil {
			fmt.Println(err)
			continue
		}

		for _, item := range items {
			svc.repo.Create(item)
		}
		time.Sleep(6 * time.Hour)
	}
}
