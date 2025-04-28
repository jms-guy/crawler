package main

import (
	"log"
	"net/url"
	"os"
	"strconv"
	"sync"
)

type config struct {
	//The list of pages crawled, number of times page was seen, and max number of pages to visit
	pages				map[string]int
	externalLinks		map[string]int
	errorPages			[]errorPages
	maxPages			int
	//Base URL given to crawler
	baseURL				*url.URL
	//Concurrency safety and control variables
	mu					*sync.Mutex
	concurrencyControl	chan struct{}
	wg					*sync.WaitGroup
}

//Arbitary max number of pages to set if user inputs 0 as amount
var maxNumofPages int = 10000

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

	//If user has entered a max number of pages of 0
	//Sets max to var maxNumofPages
	if maxPageNumberInt == 0 {
		maxPageNumberInt = maxNumofPages
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
		externalLinks: 		make(map[string]int),
		errorPages: 		make([]errorPages, 0),
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

	//Print page results to stdout
	cfg.printReport(cfg.pages, baseUrlInput)
}