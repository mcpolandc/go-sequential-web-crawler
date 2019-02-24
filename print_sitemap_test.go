package main

import "testing"

func TestMapPrettyPrint(t *testing.T) {
	PrintSitemap(mockMap, "https://golang.org/")
}

// mocks & test data
var mockMap = map[string]interface{}{
	"https://golang.org/": []string{
		"https://golang.org/pkg/",
		"https://golang.org/cmd/",
	},
	"https://golang.org/pkg/": []string{
		"https://golang.org/",
		"https://golang.org/cmd/",
		"https://golang.org/pkg/fmt/",
		"https://golang.org/pkg/os/",
	},
	"https://golang.org/pkg/fmt/": []string{
		"https://golang.org/",
		"https://golang.org/pkg/",
	},
	"https://golang.org/pkg/os/": []string{
		"https://golang.org/",
		"https://golang.org/pkg/",
	},
}
