package rss_parser

import (
	"context"
)

// Result is a struct that holds the result of a job
type Result struct {
	Items []RssItem
	Err   error
}

// Job is a struct that holds the job to be executed
type Job struct {
	Url string
}

// execute executes the job
func (j Job) execute(ctx context.Context) Result {
	feed, err := parseURLWithContext(ctx, j.Url)
	if err != nil {
		return Result{
			Err: err,
		}
	}
	rssItems := make([]RssItem, 0)
	for _, item := range feed.Items {
		rssItems = append(rssItems, RssItem{
			Title:       item.Title,
			Source:      feed.Title,
			SourceURL:   j.Url,
			Link:        item.Link,
			PublishDate: item.Date,
			Description: item.Description,
		})
	}
	return Result{
		Items: rssItems,
	}
}
