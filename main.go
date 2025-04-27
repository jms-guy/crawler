package main

import (
	"log"
	"net/url"
	"os"
	"strconv"
	"sync"
)

//Will crawl a given website, avoiding links that take it to a new domain.

type config struct {
	//The list of pages crawled, number of times page was seen, and max number of pages to visit
	pages				map[string]int
	maxPages			int
	//Base URL given to crawler
	baseURL				*url.URL
	//Concurrency safety and control variables
	mu					*sync.Mutex
	concurrencyControl	chan struct{}
	wg					*sync.WaitGroup
}

func main() {
	//User arguments
	//[program name] [website] [max concurrency value] [max number of pages to crawl]
	args := os.Args

	if len(args) < 4 {
		log.Println("Command syntax: [program name] [website] [max concurrency value] [max number of pages to crawl]")
		os.Exit(1)
	} else if len(args) > 4 {
		log.Println("Too many arguments provided")
		os.Exit(1)
	}
	baseUrlInput := args[1]
	maxConcurency := args[2]
	maxPageNumber := args[3]

	//Convert given string value inputs to ints
	maxConcurencyInt, err := strconv.Atoi(maxConcurency)
	if err != nil {
		log.Printf("Error converting given concurrency value to int: %s", err)
		os.Exit(1)
	}

	maxPageNumberInt, err := strconv.Atoi(maxPageNumber)
	if err != nil {
		log.Printf("Error converting given max number of pages to int: %s", err)
		os.Exit(1)
	}

	//Parse URL argument for url.URL struct
	baseUrl, err := url.Parse(baseUrlInput)
	if err != nil {
		log.Printf("Error parsing given URL: %s", err)
		os.Exit(1)
	}
	
	//Initialize config
	cfg := config{
		pages: 				make(map[string]int),
		maxPages: 			maxPageNumberInt,
		baseURL: 			baseUrl,
		mu: 				&sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurencyInt), //Number of goroutines crawler can run at once
		wg: 				&sync.WaitGroup{},
	}

	//Normalize the base URL
	base, err := normalizeURL(baseUrl)
	if err != nil {
		log.Printf("Error normalizing given URL: %s", err)
		os.Exit(1)
	} 

	//Start crawler and concurrency WaitGroup
	log.Printf("Starting crawl of website: %s", baseUrl)
	cfg.wg.Add(1)
	cfg.crawlPage(base)
	cfg.wg.Wait()

	if len(cfg.pages) == 0 {
		log.Printf("Empty map, something went wrong.")
		os.Exit(1)
	}

	log.Printf("%d pages visisted.", len(cfg.pages))
	//Print page results to stdout
	for key, val := range cfg.pages {
		log.Printf("\rPage: %s seen %d time(s).", key, val)
	}
}