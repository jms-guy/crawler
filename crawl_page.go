package main

import (
	"log"
	"net/url"
)

func (cfg *config) crawlPage(rawCurrentURL string) {
	//Check if max number of pages to visit has been reached
	cfg.mu.Lock()
	if len(cfg.pages) >= cfg.maxPages {
		cfg.mu.Unlock()
		cfg.wg.Done()
		return
	}
	cfg.mu.Unlock()

	//Send struct into channel, defering receiving it and WaitGroup.Done
	cfg.concurrencyControl <- struct{}{}
	defer cfg.wg.Done()
	defer func() {<-cfg.concurrencyControl}()

	//Parse current URL into url.URL struct
	currentURLData, err := url.Parse(rawCurrentURL)
	if err != nil {
		log.Printf("Error parsing currentURL: %s", err)
		return
	}

	//If current URL leads off-domain, we don't want to crawl it
	if currentURLData.Host != cfg.baseURL.Host {
		return
	}

	//Normalize the current URL
	currentURLNormalized, err := normalizeURL(currentURLData)
	if err != nil {
		log.Printf("Error normalizing currentURL: %s", err)
		return
	}

	//Check config map to see if URL has been visited
	isFirst := cfg.addPageVisit(currentURLNormalized)
	if !isFirst {
		return
	}
	
	//Check again if max pages has been reached
	cfg.mu.Lock()
	if len(cfg.pages) >= cfg.maxPages {
		cfg.mu.Unlock()
		return
	}
	cfg.mu.Unlock()

	//Get HTML from current URL
	log.Printf("\rGetting HTML from URL: %s", currentURLNormalized)
	htmlResult, err := getHTML(currentURLNormalized)
	if err != nil {
		log.Printf("Error getting html from current url: %s", err)
		return
	}
	log.Printf("\rGrabbed HTML successfully from URL: %s", currentURLNormalized)

	//Get URLs from grabbed HTML 
	nextURLS, err := getURLsFromHTML(htmlResult, cfg.baseURL)
	if err != nil {
		log.Printf("Error getting URLs from HTML data from: %s", currentURLNormalized)
		return
	}

	//Iterate over grabbed URLs, beginning goroutines for recursive calls
	for _, link := range nextURLS {
		cfg.wg.Add(1)
		go cfg.crawlPage(link)
	}
}

//If page hasn't been visited, add URL to map with value of 1.
//If page has been visited before, increment URL value and return.
func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	//Lock mutex for concurrency safety
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	_, ok := cfg.pages[normalizedURL]
	if ok {
		cfg.pages[normalizedURL]++
		return false
	}
	cfg.pages[normalizedURL] = 1
	return true
}