package main

import (
	"log"
	"os"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		log.Println("no website provided")
		os.Exit(1)
	} else if len(args) > 2 {
		log.Println("too many arguments provided")
		os.Exit(1)
	}

	baseUrl := args[1]

	log.Printf("starting crawl of: %s", baseUrl)
	htmlResult, err := getHTML(baseUrl)
	if err != nil {
		log.Printf("error getting html from %s: %s", baseUrl, err)
		os.Exit(1)
	}

	log.Println(htmlResult)
}