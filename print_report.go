package main

import (
	"log"
	"sort"
)

//Struct to hold results of crawled link maps
type webpages struct{
	page		string
	visited		int
}

func (cfg *config) printReport(pages map[string]int, baseURL string) {
	//Results for found links on-domain
	log.Println("\r=============================")
	log.Printf("\r  REPORT for %s", baseURL)
	log.Println("\r=============================")

	srtdPages := sortPages(pages)

	for _, page := range srtdPages{
		log.Printf("\rFound %d internal link(s) to %s", page.visited, page.page)
	}

	//Results for external links found
	log.Println("\r=============================")
	log.Println("\r  External links found")
	log.Println("\r=============================")

	srtdExtLinks := sortPages(cfg.externalLinks)

	for _, page := range srtdExtLinks {
		log.Printf("\rFound %d external link(s) to %s", page.visited, page.page)
	}

	//Results for links that produced errors on crawling
	log.Println("\r=============================")
	log.Println("\r  Broken/Error'd links encountered")
	log.Println("\r=============================")

	cfg.sortErrorPages()

	for _, errPage := range cfg.errorPages{
		log.Printf("\rFound %d link(s) to %s: %s", errPage.visited, errPage.page, errPage.reason)
	}
}

//Sorts map of links into struct by highest count, then alphabetically
func sortPages(pages  map[string]int) []webpages {
	var sorted []webpages

	//Create new struct in slice for each URL
	for page, num := range pages{
		newPage := webpages{
			visited: num,
			page: page,
		}
		sorted = append(sorted, newPage)
	}
	
	sort.Slice(sorted, func (i, j int) bool {
		if sorted[i].visited != sorted[j].visited {
			return sorted[i].visited > sorted[j].visited
		}
		return sorted[i].page < sorted[j].page
	})
	return sorted
}