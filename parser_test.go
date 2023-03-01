package rss_parser_test

import (
	rssparser "github.com/simplexpage/rss-parser"
	"testing"
)

func TestParseURLs(t *testing.T) {
	rssUrls := []string{
		"https://www.reddit.com/r/golang/.rss",
		"https://www.reddit.com/r/golang/new/.rss",
		"https://www.reddit.com/r/golang/top/.rss",
		"https://www.reddit.com/r/golang/comments/.rss",
		"https://www.reddit.com/r/golang/controversial/.rss",
	}
	rssUrlsParser := rssparser.NewParser()
	rssItems, err := rssUrlsParser.ParseURLs(rssUrls)
	if err != nil {
		t.Error(err)
	}
	if len(rssItems) == 0 {
		t.Error("No rss items found")
	}
	t.Log(rssItems)
}
