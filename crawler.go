package main

import (
	"log"
	"net/http"
	"sync"

	"golang.org/x/net/html"
)

// ICrawler - interface describing behaviour of crawler
type ICrawler interface {
	Crawl(url string, wg *sync.WaitGroup) error
	Fetch(url string) ([]string, error)
}

// Crawler - an implementation of the Crawler interface
type Crawler struct {
	Domain      string
	CrawledUrls *ThreadSafeMap // using a map for performance reasons (O(1))
}

// Crawl - recursively crawl given urls
func (crawlerPtr *Crawler) Crawl(url string, wg *sync.WaitGroup) error {

	wg.Add(1)
	defer wg.Done()

	log.Printf("Crawling %s\n", url)

	webCrawler := *crawlerPtr
	fullURL := PrependDomain(url, webCrawler.Domain)

	// Get slice of urls from html
	foundsUrls, err := webCrawler.Fetch(fullURL)
	if err != nil {
		return err
	}

	// Add url to list of already crawled
	webCrawler.CrawledUrls.Set(fullURL, foundsUrls)

	// Iterate over slice of found urls
	for _, link := range foundsUrls {
		// Check url has not already been crawled

		if webCrawler.CrawledUrls.Get(link) != nil {
			go webCrawler.Crawl(link, wg)
		}
	}

	return err
}

// Fetch - return a slice of urls found in response body
func (crawlerPtr *Crawler) Fetch(fullURL string) (urls []string, err error) {

	// Fetch html from url
	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, err
	}

	// Use tokenizer to get html elements from resp
	tokenizer := html.NewTokenizer(resp.Body)

	// Iterate over tokens
	for {

		tokenType := tokenizer.Next()

		if tokenType == html.StartTagToken {

			token := tokenizer.Token()

			// Find anchor type elements
			if token.Data == "a" {
				// Iterate over anchors to eligible links
				for _, attr := range token.Attr {
					if isEligibleLink(attr) {
						urls = append(urls, PrependDomain(attr.Val, (*crawlerPtr).Domain))
					}
				}
			}
		}

		if tokenType == html.ErrorToken {
			// Return urls slice
			return urls, nil
		}
	}
}

func isEligibleLink(attr html.Attribute) bool {
	return IsHref(attr) &&
		IsInternalLink(attr.Val) &&
		!IsSamePageLink(attr.Val) &&
		!IsPhoneLink(attr.Val)
}
