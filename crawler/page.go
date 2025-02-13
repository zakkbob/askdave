package main

import (
	"sync"
)

type page struct {
	url           string
	pageTitle     string
	ogTitle       string
	ogDescription string
	ogSiteName    string
	links         []url
}

type safePageSlice struct {
	mutex sync.Mutex
	slice []page
}

func (p *page) addLink(u url) {
	p.links = append(p.links, u)
}
