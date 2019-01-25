package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

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
				//    2B. iterate over anchors to find hrefs
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						//    2C. ensure href is an internal link
						if !strings.HasPrefix(attr.Val, "http") {
							//    2D. save link to urls slice
							urls = append(urls, prependDomain(attr.Val, domain))
						}
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

// util
func prependDomain(route string, domain string) (url string) {
	if strings.HasPrefix(route, "http") {
		url = route
	} else {
		url = fmt.Sprintf("%s%s", domain, route)
	}
	return url
}

// Crawler - interface describing behaviour of crawler
type Crawler interface {
	Crawl(url string) error
}

// CrawlerImpl - an implementation of the Crawler interface
type CrawlerImpl struct {
	domain      string
	data        map[string]*result
	crawledUrls []string
}

type result struct {
	links []string
}

// Crawl - recursively crawl given urls
func (webCrawler *CrawlerImpl) Crawl(url string) error {

	fullURL := prependDomain(url, webCrawler.domain)

	// 1. fetch html from url
	resp, err := http.Get(fullURL)
	if err != nil {
		return err
	}

	// 2. get slice of urls from html
	foundsUrls, err := Fetch(resp, (*webCrawler).domain)
	if err != nil {
		return err
	}

	// 3. add url to list of already crawled
	(*webCrawler).crawledUrls = append((*webCrawler).crawledUrls, fullURL)

	// 4. add urls to map with this url as key
	(*webCrawler).data[fullURL] = &result{links: foundsUrls}

	// 5. iterate over slice of found urls
	for _, link := range foundsUrls {
		//    5A. check url has not already been crawled

		if !Contains((*webCrawler).crawledUrls, link) {
			(*webCrawler).Crawl(link)
		}
	}

	fmt.Printf("(*webCrawler).crawledUrls: %v", (*webCrawler).crawledUrls)
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
		domain: os.Args[1],
		data:   make(map[string]*result),
	}

	crawler.Crawl(domain)
}
