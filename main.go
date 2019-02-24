package main

import (
	"fmt"
	"log"
	"os"
	"sync"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Program did not receive any args")
		os.Exit(1)
	}

	domain := os.Args[1]

	crawler := Crawler{
		Domain:      os.Args[1],
		CrawledUrls: &ThreadSafeMap{items: make(map[string]interface{})},
	}

	log.Printf("starting to crawl %v", domain)

	var wg sync.WaitGroup

	wg.Add(1)
	defer wg.Done()

	crawler.Crawl(domain, &wg)

	wg.Wait()

	log.Printf("Finished crawling. Crawled %d pages", len(crawler.CrawledUrls.items))
}
