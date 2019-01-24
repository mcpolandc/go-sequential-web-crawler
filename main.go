package main

import (
	"fmt"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

type fakeResult struct {
	body string
	urls []string
}

type fakeFetcher map[string]*fakeResult

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (urls []string, err error)
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) {
	// TODO: Fetch URLs in parallel.
	// TODO: Don't fetch the same URL twice.
	// This implementation doesn't do either:
	if depth <= 0 {
		return
	}

	urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("found: %s\n", url)
	for _, u := range urls {
		Crawl(u, depth-1, fetcher)
	}
	return
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Program did not receive an argument")
		os.Exit(1)
	}

	domain := os.Args[1]

	fetcher := new(fakeFetcher)
	Crawl(domain, 4, fetcher)
}

// fakeFetcher is Fetcher that returns canned results.

func (urlStore *fakeFetcher) Fetch(url string) ([]string, error) {

	page, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	tokenizer := html.NewTokenizer(page.Body)

	for {
		tokenType := tokenizer.Next()

		if tokenType == html.ErrorToken {
			return nil, nil
		}

		if tokenType == html.StartTagToken {
			token := tokenizer.Token()

			if token.Data == "a" {
				for _, attr := range token.Attr {
					if attr.Key == "href" { // also check if internal
						// *urlStore[url] = &{something like this!}
					}
				}
			}
		}
	}
}

// fetcher is a populated fakeFetcher.
var fake = fakeFetcher{
	"httfetcherps://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}
