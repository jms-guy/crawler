package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func getHTML(rawURL string) (string, error) {
	res, err := http.Get(rawURL)
	if err != nil {
		log.Printf("Error sending http request to %s: %s", rawURL, err)
		return "", err
	}

	if res.StatusCode >= 400 {
		return "", fmt.Errorf("response status code 400+")
	}

	header := res.Header.Get("content-type")
	if header != "text/html" {
		return "", fmt.Errorf("content type not html format")
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response data: %s", err)
	}

	htmlString := string(bytes)

	return htmlString, nil
}