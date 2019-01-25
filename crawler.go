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
}

// Crawler - an implementation of the Crawler interface
type Crawler struct {
	Domain      string
	Data        *ThreadSafeMap
	CrawledUrls *ThreadSafeMap // using a map for performance reasons (O(1))
}

// Crawl - recursively crawl given urls
func (crawlerPtr *Crawler) Crawl(url string, wg *sync.WaitGroup) error {

	defer wg.Done()

	log.Printf("Crawling %s\n", url)

	webCrawler := *crawlerPtr
	fullURL := PrependDomain(url, webCrawler.Domain)

	// 1. fetch html from url
	resp, err := http.Get(fullURL)
	if err != nil {
		return err
	}

	// 2. get slice of urls from html
	foundsUrls, err := Fetch(resp, webCrawler.Domain)
	if err != nil {
		return err
	}

	// 3. add url to list of already crawled
	webCrawler.CrawledUrls.Set(fullURL, true)

	// 4. add urls to map with this url as key
	webCrawler.Data.Set(fullURL, foundsUrls)

	// 5. iterate over slice of found urls
	for _, link := range foundsUrls {
		//    5A. check url has not already been crawled

		if webCrawler.CrawledUrls.Get(link) != nil {
			go webCrawler.Crawl(link, wg)
		}
	}

	return err
}

// TODO - Move to own file
// Fetch - return a slice of urls found in response body
func Fetch(resp *http.Response, domain string) (urls []string, err error) {

	// 1. use tokenizer to get html elements from resp
	tokenizer := html.NewTokenizer(resp.Body)

	// 2. iterate over tokens
	for {

		tokenType := tokenizer.Next()

		if tokenType == html.StartTagToken {

			token := tokenizer.Token()

			//    2A. find anchor type elements
			if token.Data == "a" {
				//    2B. iterate over anchors to eligible links
				for _, attr := range token.Attr {
					if isEligibleLink(attr) {
						urls = append(urls, PrependDomain(attr.Val, domain))
					}
				}
			}
		}

		if tokenType == html.ErrorToken {
			//    2E. return urls slice
			return urls, nil
		}
	}
}

func isEligibleLink(attr html.Attribute) bool {
	return IsHref(attr) &&
		IsInternalLink(attr.Val) &&
		!IsSamePageLink(attr.Val) &&
		!IsScriptLink(attr.Val)
}
