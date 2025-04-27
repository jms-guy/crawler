package main

import (
	"log"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	reader := strings.NewReader(htmlBody)
	htmlNodes, err := html.Parse(reader)
	if err != nil {
		log.Printf("Error parsing html reader: %s", err)
		return nil, err
	}

	baseUrl, err := url.Parse(rawBaseURL)
	if err != nil {
		log.Printf("Error parsing baseURL %s: %s", rawBaseURL, err)
		return nil, err
	}

	var urls []string

	for node := range htmlNodes.Descendants() {
		if node.Type == html.ElementNode && node.Data == "a" {
			for _, a := range node.Attr {
				if a.Key == "href" {
					parsedLink, err := url.Parse(a.Val)
					if err != nil {
						continue
					}
					resolvedUrl := baseUrl.ResolveReference(parsedLink)
					urls = append(urls, resolvedUrl.String())
				}
			}
		}
	}
	return urls, nil
}