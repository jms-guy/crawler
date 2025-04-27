package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func getHTML(rawURL string) (string, error) {
	//Make http request to given URL
	res, err := http.Get(rawURL)
	if err != nil {
		log.Printf("Error sending http request to %s: %s", rawURL, err)
		return "", err
	}

	//Check response status code
	if res.StatusCode >= 400 {
		return "", fmt.Errorf("response status code 400+")
	}

	//Check response header - only want HTML
	header := res.Header.Get("content-type")
	if !strings.Contains(header, "text/html") {
		return "", fmt.Errorf("content type not html format")
	}

	//Read response body
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response data: %s", err)
	}

	//Stringify HTML response data
	htmlString := string(bytes)

	return htmlString, nil
}