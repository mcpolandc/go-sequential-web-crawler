package main

import "testing"

func TestCrawl(t *testing.T) {
	body, urls, _ := fetcher.Fetch("https://golang.org/")

	if body != "The Go Programming Language" {
		t.Fail()
	}

	if len(urls) != 2 {
		t.Fail()
	}
}
