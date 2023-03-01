package rss_parser

import (
	"time"
)

type parser struct{}

func NewParser() *parser {
	return &parser{}
}

func (f *parser) ParseURLs(urls []string) ([]RssItem, error) {
	var rssItems []RssItem
	for _, url := range urls {
		rssItems = append(rssItems, RssItem{
			Title:       "title",
			Source:      "source",
			SourceURL:   url,
			Link:        "link",
			PublishDate: time.Now(),
			Description: "description",
		})
	}
	return rssItems, nil
}

type RssItem struct {
	Title       string
	Source      string
	SourceURL   string
	Link        string
	PublishDate time.Time
	Description string
}
