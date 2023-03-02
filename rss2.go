package rss_parser

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"io"
	"io/ioutil"
	"strings"
	"time"
)

func parseRSS2(data []byte) (*Feed, error) {
	rss2Feed := rss2Feed{}
	p := xml.NewDecoder(bytes.NewReader(data))
	p.CharsetReader = charsetReader
	err := p.Decode(&rss2Feed)
	if err != nil {
		return nil, err
	}
	if rss2Feed.Channel == nil {
		return nil, fmt.Errorf("no channel found in %q", string(data))
	}

	channel := rss2Feed.Channel

	feed := new(Feed)
	feed.Title = channel.Title
	feed.Link = channel.Link
	feed.Items = make([]*Item, 0, len(channel.Items))
	for _, item := range channel.Items {
		nextItem := new(Item)
		nextItem.Title = item.Title
		nextItem.Description = item.Description
		nextItem.Link = item.Link
		nextItem.Date = parseDate(item.PubDate)
		feed.Items = append(feed.Items, nextItem)
	}
	return feed, nil
}

type rss2Feed struct {
	XMLName xml.Name     `xml:"rss"`
	Channel *rss2Channel `xml:"channel"`
}

type rss2Channel struct {
	XMLName xml.Name   `xml:"channel"`
	Title   string     `xml:"title"`
	Link    string     `xml:"link"`
	Items   []rss2Item `xml:"item"`
}

type rss2Item struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`
	Description string   `xml:"description"`
	Link        string   `xml:"link"`
	PubDate     string   `xml:"pubDate"`
	Date        string   `xml:"date"`
}

func parseDate(date string) time.Time {
	dateNew, err := time.Parse(time.RFC1123Z, date)
	if err != nil {
		dateNew = time.Now()
	}
	return dateNew
}

func charsetReader(charset string, input io.Reader) (io.Reader, error) {
	switch {
	case isCharsetUTF8(charset):
		return input, nil
	case isCharsetWindows1251(charset):
		return newUT8CharsetFromWindows1251(input), nil
	}
	// TODO: implement other charsets
	return nil, errors.New("CharsetReader: unexpected charset: " + charset)
}

func isCharset(charset string, names []string) bool {
	charset = strings.ToLower(charset)
	for _, n := range names {
		if charset == strings.ToLower(n) {
			return true
		}
	}
	return false
}

func isCharsetUTF8(charset string) bool {
	names := []string{
		"UTF-8",
		"",
	}
	return isCharset(charset, names)
}

func isCharsetWindows1251(charset string) bool {
	names := []string{
		"windows-1251",
	}
	return isCharset(charset, names)
}

func newUT8CharsetFromWindows1251(input io.Reader) io.Reader {
	decoder := charmap.Windows1251.NewDecoder()
	reader := decoder.Reader(input)
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(b)
}
