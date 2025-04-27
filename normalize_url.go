package main

import (
	"net/url"
	"strings"
)

//Take url.URL struct and normalize URL data into
// [scheme][host]/[path] format
func normalizeURL(rawURL *url.URL) (string, error) {
	rawURL.Path = strings.TrimRight(rawURL.Path, "/")
	normalizedUrl := "http://" + rawURL.Host + rawURL.Path
	return normalizedUrl, nil
}