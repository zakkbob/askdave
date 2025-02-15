//-------------------------------//
// Works through all crawl tasks //
//-------------------------------//

package main

import (
	"fmt"
	"time"
)

var crawledPages = safePageSlice{slice: make([]page, 0)}

var robotsCrawlTasks = safeUrlSlice{slice: make([]url, 0)}
var sitemapCrawlTasks = safeUrlSlice{slice: make([]url, 0)}
var pageCrawlTasks = safeUrlSlice{slice: make([]url, 0)}

// var robotsCrawlTasks = safeUrlSlice{slice: make([]url, 0)}
// var sitemapCrawlTasks = safeUrlSlice{slice: make([]url, 0)}
// var pageCrawlTasks = safeUrlSlice{slice: make([]url, 0)}

func completeNextTask() {
	//Do robots crawl if present
	robotsCrawlTasks.mu.Lock()
	if len(robotsCrawlTasks.slice) > 0 {
		u := robotsCrawlTasks.slice[0]
		robotsCrawlTasks.slice = robotsCrawlTasks.slice[1:len(robotsCrawlTasks.slice)]
		robotsCrawlTasks.mu.Unlock()
		completeRobotsCrawl(u)
		return
	}
	robotsCrawlTasks.mu.Unlock()

	//Do sitemap crawl if present
	sitemapCrawlTasks.mu.Lock()
	if len(sitemapCrawlTasks.slice) > 0 {
		u := sitemapCrawlTasks.slice[0]
		sitemapCrawlTasks.slice = sitemapCrawlTasks.slice[1:len(sitemapCrawlTasks.slice)]
		sitemapCrawlTasks.mu.Unlock()
		completeSitemapCrawl(u)
		return
	}
	sitemapCrawlTasks.mu.Unlock()

	//Do page crawl if present
	pageCrawlTasks.mu.Lock()
	if len(pageCrawlTasks.slice) > 0 {
		u := pageCrawlTasks.slice[0]
		pageCrawlTasks.slice = pageCrawlTasks.slice[1:len(pageCrawlTasks.slice)]
		pageCrawlTasks.mu.Unlock()
		completePageCrawl(u)
		return
	}
	pageCrawlTasks.mu.Unlock()
}

func completeRobotsCrawl(u url) {
	time.Sleep(1000 * time.Millisecond)
	fmt.Printf("Crawled robots.txt '%s'\n", u.String())
}

func completeSitemapCrawl(u url) {
	time.Sleep(1000 * time.Millisecond)
	fmt.Printf("Crawled sitemap '%s'\n", u.String())
}

func completePageCrawl(u url) {
	time.Sleep(1000 * time.Millisecond)
	fmt.Printf("Crawled page '%s'\n", u.String())
}
