package main

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

// IsHref - takes html attribute and checks if is of type href
func IsHref(attr html.Attribute) bool {
	return attr.Key == "href"
}

// IsSamePageLink - checks if link points somewhere on the same page
func IsSamePageLink(link string) bool {
	return strings.HasPrefix(link, "#")
}

// IsInternalLink - checks if internal link, assume internal links are prefixed with '/'
func IsInternalLink(link string) bool {
	return strings.HasPrefix(link, "/")
}

// IsScriptLink - checks if cdn-cgi script links are present
func IsScriptLink(link string) bool {
	return strings.HasPrefix(link, "/cdn-cgi")
}

// PrependDomain - takes a route and adds domain if needed
func PrependDomain(route string, domain string) (url string) {
	if strings.HasPrefix(route, "http") {
		url = route
	} else {
		url = fmt.Sprintf("%s%s", domain, route)
	}
	return url
}
