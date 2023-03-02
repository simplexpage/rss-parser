package rss_parser_test

import (
	rssparser "github.com/simplexpage/rss-parser"
	"testing"
)

func TestParseURLs(t *testing.T) {
	rssUrls := []string{
		"https://tsn.ua/rss/full.rss",
		"https://www.pravda.com.ua/rus/rss/",
	}
	rssUrlsParser := rssparser.NewRssParser()
	rssItems, err := rssUrlsParser.ParseURLs(rssUrls)
	if err != nil {
		t.Error(err)
	}
	if len(rssItems) == 0 {
		t.Error("No rss items found")
	}
	t.Log(rssItems)
}
