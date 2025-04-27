package main

import (
	"log"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func getURLsFromHTML(htmlBody string, rawBaseURL *url.URL) ([]string, error) {
	//Create io.Reader from HTML body string
	reader := strings.NewReader(htmlBody)

	//Parse io.Reader into HTML nodes 
	htmlNodes, err := html.Parse(reader)
	if err != nil {
		log.Printf("Error parsing html reader: %s", err)
		return nil, err
	}

	//Initialize return slice of URLS
	var urls []string

	//Iterate over all HTML nodes, check for link tags <a href=
	//If link tag is found, resolve the URL based on initial
	//base URL and add to return slice
	for node := range htmlNodes.Descendants() {
		if node.Type == html.ElementNode && node.Data == "a" {
			for _, a := range node.Attr {
				if a.Key == "href" {
					parsedLink, err := url.Parse(a.Val)
					if err != nil {
						continue
					}
					resolvedUrl := rawBaseURL.ResolveReference(parsedLink)
					urls = append(urls, resolvedUrl.String())
				}
			}
		}
	}
	return urls, nil
}