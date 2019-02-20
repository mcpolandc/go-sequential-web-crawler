package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

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
		!IsScriptLink(attr.Val) &&
		!IsPhoneLink(attr.Val)
}

// Crawler - interface describing behaviour of crawler
type Crawler interface {
	Crawl(url string) error
}

// CrawlerImpl - an implementation of the Crawler interface
type CrawlerImpl struct {
	domain      string
	data        map[string]*CrawlResult
	crawledUrls map[string]bool // using a map for performance reasons (O(1))
}

// CrawlResult - store links associated with a url
type CrawlResult struct {
	links []string
}

// Crawl - recursively crawl given urls
func (crawlerPtr *CrawlerImpl) Crawl(url string) error {

	webCrawler := *crawlerPtr
	fullURL := PrependDomain(url, webCrawler.domain)

	// 1. fetch html from url
	resp, err := http.Get(fullURL)
	if err != nil {
		return err
	}

	// 2. get slice of urls from html
	foundsUrls, err := Fetch(resp, webCrawler.domain)
	if err != nil {
		return err
	}

	// 3. add url to list of already crawled
	webCrawler.crawledUrls[fullURL] = true

	// 4. add urls to map with this url as key
	webCrawler.data[fullURL] = &CrawlResult{links: foundsUrls}

	// 5. iterate over slice of found urls
	for _, link := range foundsUrls {
		//    5A. check url has not already been crawled

		if !webCrawler.crawledUrls[link] {
			webCrawler.Crawl(link)
		}
	}

	return err
}

// Contains - checks is slice `a` contains string `x`
func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Program did not receive any args")
		os.Exit(1)
	}

	domain := os.Args[1]

	crawler := CrawlerImpl{
		domain:      os.Args[1],
		data:        make(map[string]*CrawlResult),
		crawledUrls: make(map[string]bool),
	}

	log.Printf("starting to crawl %v", domain)

	crawler.Crawl(domain)

	PrintSitemap(crawler.data, domain)

	log.Printf("Finished crawling. Crawled %d pages", len(crawler.data))
}
