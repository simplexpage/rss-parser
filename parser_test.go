package rss_parser_test

import (
	"context"
	"fmt"
	rssparser "github.com/simplexpage/rss-parser"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

// Test parseRSS
func TestParseURLsSuccess(t *testing.T) {
	// Test cases
	var feedsTests = []struct {
		file  string
		items []rssparser.RssItem
	}{
		{"rss2.xml", []rssparser.RssItem{
			{
				"Example item 1",
				"RSS2 Title",
				"http:://example.com/rss",
				"http://example.com/1",
				time.Date(2023, 01, 01, 00, 01, 00, 0, time.UTC),
				"Here is some text 1.",
			},
			{
				"Example item 2",
				"RSS2 Title",
				"http:://example.com/rss",
				"http://example.com/2",
				time.Date(2023, 01, 01, 00, 01, 00, 0, time.UTC),
				"Here is some text 2.",
			},
		},
		},
	}

	// Create mock servers
	var feedsUrls []string
	for i, test := range feedsTests {
		path := fmt.Sprintf("testdata/%s", test.file)
		f, _ := ioutil.ReadFile(path)
		server, _ := mockServerResponse(200, string(f), 0)
		feedsUrls = append(feedsUrls, server.URL)
		for j, _ := range test.items {
			feedsTests[i].items[j].SourceURL = server.URL
		}
	}

	feed, err := rssparser.ParseURLs(context.Background(), feedsUrls)

	assert.NotNil(t, feed)
	assert.Nil(t, err)
	assert.Equal(t, feedsTests[0].items, feed)
}

func TestParseURLsFailure(t *testing.T) {
	server, _ := mockServerResponse(404, "", 0)
	feed, err := rssparser.ParseURLs(context.Background(), []string{server.URL})

	assert.NotNil(t, err)
	assert.IsType(t, rssparser.HTTPError{}, err)
	assert.Nil(t, feed)
}

func TestParser_ParseURLWithContext(t *testing.T) {
	server, _ := mockServerResponse(404, "", 1*time.Minute)
	ctxTime, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := rssparser.ParseURLs(ctxTime, []string{server.URL})
	assert.True(t, strings.Contains(err.Error(), "context deadline exceeded"))
}

// Test Helpers
func mockServerResponse(code int, body string, delay time.Duration) (*httptest.Server, *http.Client) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delay)
		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/xml")
		io.WriteString(w, body)
	}))

	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}

	client := &http.Client{Transport: transport}
	return server, client
}

func ExampleRssParser_ParseURLs() {
	rssUrls := []string{
		"https://tsn.ua/rss/full.rss",
		"https://www.pravda.com.ua/rus/rss/",
	}

	ctxTime, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rssItems, err := rssparser.ParseURLs(ctxTime, rssUrls)
	if err != nil {
		fmt.Println(err)
	}
	for _, item := range rssItems {
		fmt.Println(item.Title)
	}
}
