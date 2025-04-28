package main

import (
	"sort"
)

//Struct to hold links that produce errors when crawling
type errorPages struct{
	page		string
	visited		int
	reason		string
}

func (cfg *config) handleErrorPages(pageURL string, err error) {
	//If link is already present in main []errorPages, increments it
	for i, ep := range cfg.errorPages {
		if ep.page == pageURL {
			cfg.errorPages[i].visited++
			return
		}
	}
	//Else adds link to it
	newPage := errorPages{
		page: pageURL,
		visited: 1,
		reason: err.Error(),
	}
	cfg.errorPages = append(cfg.errorPages, newPage)
}

//Sorts the links in struct, highest count first, then alphabetically
func (cfg *config) sortErrorPages() {
	sort.Slice(cfg.errorPages, func (i, j int) bool {
		if cfg.errorPages[i].visited != cfg.errorPages[j].visited {
			return cfg.errorPages[i].visited > cfg.errorPages[j].visited
		}
		return cfg.errorPages[i].page < cfg.errorPages[j].page
	})
}