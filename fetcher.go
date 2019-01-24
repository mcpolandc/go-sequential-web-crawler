// package main

// FetchResult is returned for a given url
type FetchResult struct {
	body string
	urls []string
}

// Fetcher interface
type Fetcher interface {
	Fetch(url string) (body string, urls []string, err error)
}

// func (f Fetcher) Fetch(url string) (string, []string, error) {

// }
