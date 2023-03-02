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

// parseRSS2 parses rss 2.0 feed from bytes
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

// rss2Feed is a struct for rss 2.0 feed
type rss2Feed struct {
	XMLName xml.Name     `xml:"rss"`
	Channel *rss2Channel `xml:"channel"`
}

// rss2Channel is a struct for rss 2.0 channel
type rss2Channel struct {
	XMLName xml.Name   `xml:"channel"`
	Title   string     `xml:"title"`
	Link    string     `xml:"link"`
	Items   []rss2Item `xml:"item"`
}

// rss2Item is a struct for rss 2.0 item
type rss2Item struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`
	Description string   `xml:"description"`
	Link        string   `xml:"link"`
	PubDate     string   `xml:"pubDate"`
	Date        string   `xml:"date"`
}

// parseDate parses date from string
func parseDate(date string) time.Time {
	dateNew, err := time.Parse(time.RFC1123Z, date)
	if err != nil {
		dateNew = time.Now()
	}
	return dateNew
}

// charsetReader is a function for reading charset
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

// isCharset checks if charset is in names
func isCharset(charset string, names []string) bool {
	charset = strings.ToLower(charset)
	for _, n := range names {
		if charset == strings.ToLower(n) {
			return true
		}
	}
	return false
}

// isCharsetUTF8 checks if charset is UTF-8
func isCharsetUTF8(charset string) bool {
	names := []string{
		"UTF-8",
		"",
	}
	return isCharset(charset, names)
}

// isCharsetWindows1251 checks if charset is Windows-1251
func isCharsetWindows1251(charset string) bool {
	names := []string{
		"windows-1251",
	}
	return isCharset(charset, names)
}

// newUT8CharsetFromWindows1251 converts Windows-1251 to UTF-8
func newUT8CharsetFromWindows1251(input io.Reader) io.Reader {
	decoder := charmap.Windows1251.NewDecoder()
	reader := decoder.Reader(input)
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(b)
}
