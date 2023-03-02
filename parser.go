package rss_parser

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
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

func (f *RssParser) ParseURLs(urls []string) ([]RssItem, error) {
	var rssItems []RssItem
	for _, url := range urls {
		feed, err := f.parseURLWithContext(url, context.Background())
		if err != nil {
			return nil, err
		}
		for _, item := range feed.Items {
			rssItems = append(rssItems, RssItem{
				Title:       item.Title,
				Source:      feed.Title,
				SourceURL:   feed.Link,
				Link:        item.Link,
				PublishDate: item.Date,
				Description: item.Description,
			})
		}
	}
	return rssItems, nil
}

func (f *RssParser) parseURLWithContext(rssURL string, ctx context.Context) (feed *Feed, err error) {
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return f.parseRSS(body)
}

func (f *RssParser) parseRSS(data []byte) (feed *Feed, err error) {
	if strings.Contains(string(data), "<rss") {
		return parseRSS2(data)
	} else if strings.Contains(string(data), "xmlns=\"http://purl.org/rss/1.0/\"") {
		return parseRSS1(data)
	}
	// TODO: implement atom parser

	return nil, fmt.Errorf("unknown rss format")
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

type Feed struct {
	Title string  `xml:"title"`
	Link  string  `xml:"link"`
	Items []*Item `json:"items"`
}

type Item struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Link        string    `json:"link"`
	Date        time.Time `json:"date"`
}
