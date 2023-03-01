# RSS Reader package
Golang RSS Reader package, which can parse asynchronous multiple RSS feeds.

The package have only exportable type RssItem and method Parse that:
1. Takes an array of URLs.
2. Parses asynchronously their feed.
3. Returns an array of RssItem generated from all provided RSS posts.

Dependencies:
```bash
go get github.com/simplexpage/rss-parser
```

Example usage:
```go
package main
import rssparser "github.com/simplexpage/rss-parser"
func main() {
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
		// handle error.
	}
}
```

The output structure is pretty much as you'd expect:
```go
type RssItem struct{
    Title string
    Source string
    SourceURL string
    Link string
    PublishDate time.Time
    Description string
}
```