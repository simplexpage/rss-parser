package rss_parser

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// ParseURLs parses rss feeds from urls
func ParseURLs(ctx context.Context, urls []string) ([]RssItem, error) {

	if len(urls) == 0 {
		return nil, errors.New("no urls provided")
	}

	rssItems := make([]RssItem, 0)

	wp := NewWorkerPool(len(urls))

	go wp.AddJob(urls)

	go wp.Run(ctx)

	for {
		select {
		case r, ok := <-wp.Results():
			if !ok {
				continue
			}
			if r.Err != nil {
				return nil, r.Err
			}
			rssItems = append(rssItems, r.Items...)
		case <-wp.Done:
			return rssItems, nil
		default:
		}
	}
}

// parseURLWithContext requests rss feed from url
func parseURLWithContext(ctx context.Context, rssURL string) (feed *Feed, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rssURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	if resp != nil {
		defer func() {
			ce := resp.Body.Close()
			if ce != nil {
				err = ce
			}
		}()
	}

	if resp.StatusCode != http.StatusOK {
		return nil, HTTPError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
		}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return parseRSS(body)
}

// parseRSS parses rss feed from bytes
func parseRSS(data []byte) (feed *Feed, err error) {
	if strings.Contains(string(data), "<rss") {
		return parseRSS2(data)
	} else if strings.Contains(string(data), "xmlns=\"http://purl.org/rss/1.0/\"") {
		// TODO: implement rss 1.0 parser
		return nil, fmt.Errorf("unknown rss format")
	} else {
		// TODO: implement atom parser
		return nil, fmt.Errorf("unknown rss format")
	}
}

// HTTPError is an error type for http errors
type HTTPError struct {
	StatusCode int
	Status     string
}

// Error returns error message
func (err HTTPError) Error() string {
	return fmt.Sprintf("http error: %s", err.Status)
}

// RssItem is a struct for rss item
type RssItem struct {
	Title       string
	Source      string
	SourceURL   string
	Link        string
	PublishDate time.Time
	Description string
}

// Feed is a struct for rss feed
type Feed struct {
	Title string  `xml:"title"`
	Link  string  `xml:"link"`
	Items []*Item `json:"items"`
}

// Item is a struct for return rss items
type Item struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Link        string    `json:"link"`
	Date        time.Time `json:"date"`
}
