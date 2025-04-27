package main

import (
	"log"
	"net/url"
	"strings"
)

func normalizeURL(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		log.Printf("Error parsing url: %s : %s", rawURL, err)
		return "", err
	}

	u.Path = strings.TrimRight(u.Path, "/")
	normalizedUrl := u.Host + u.Path
	return normalizedUrl, nil
}