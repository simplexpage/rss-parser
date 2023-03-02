package rss_parser

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type RssParser struct {
	HttpClient *http.Client
	UserAgent  string
}

func NewRssParser() *RssParser {
	return &RssParser{
		UserAgent: "rss-parser",
	}
}

func (f *RssParser) ParseURLs(urls []string) (rssItems []RssItem, err error) {
	for _, url := range urls {
		rssItems = append(rssItems, RssItem{
			Title:       "title",
			Source:      "source",
			SourceURL:   url,
			Link:        "link",
			PublishDate: time.Now(),
			Description: "description",
		})
		f.parseURLWithContext(url, context.Background())
	}
	return rssItems, nil
}

func (f *RssParser) parseURLWithContext(rssURL string, ctx context.Context) (rssItems []RssItem, err error) {
	client := f.httpClient()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rssURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", f.UserAgent)

	resp, err := client.Do(req)

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

	if resp.StatusCode != 200 {
		return nil, HTTPError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
		}
	}

	return f.parse(resp.Body)
}

func (f *RssParser) parse(rss io.Reader) (rssItems []RssItem, err error) {
	return
}

func (f *RssParser) httpClient() *http.Client {
	if f.HttpClient != nil {
		return f.HttpClient
	}
	f.HttpClient = &http.Client{}
	return f.HttpClient
}

type HTTPError struct {
	StatusCode int
	Status     string
}

func (err HTTPError) Error() string {
	return fmt.Sprintf("http error: %s", err.Status)
}

type RssItem struct {
	Title       string
	Source      string
	SourceURL   string
	Link        string
	PublishDate time.Time
	Description string
}
