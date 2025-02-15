//-------------------------------//
// Works through all crawl tasks //
//-------------------------------//

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var crawledPages = safePageSlice{slice: make([]page, 0)}

var robotsCrawlTasks = safeUrlSlice{slice: make([]url, 0)}
var sitemapCrawlTasks = safeUrlSlice{slice: make([]url, 0)}
var pageCrawlTasks = safeUrlSlice{slice: make([]url, 0)}

// var robotsCrawlTasks = safeUrlSlice{slice: make([]url, 0)}
// var sitemapCrawlTasks = safeUrlSlice{slice: make([]url, 0)}
// var pageCrawlTasks = safeUrlSlice{slice: make([]url, 0)}

func loopConcurrently(t func() bool, workers int) {
	var wg sync.WaitGroup

	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for t() {
			}
		}()
	}

	wg.Wait()
}

func completeNextRobotsTask() bool {
	robotsCrawlTasks.mu.Lock()
	if len(robotsCrawlTasks.slice) > 0 {
		u := robotsCrawlTasks.slice[0]
		robotsCrawlTasks.slice = robotsCrawlTasks.slice[1:len(robotsCrawlTasks.slice)]
		robotsCrawlTasks.mu.Unlock()
		completeRobotsCrawl(u)
		return true
	}
	robotsCrawlTasks.mu.Unlock()
	return false
}

func completeNextSitemapTask() bool {
	sitemapCrawlTasks.mu.Lock()
	if len(sitemapCrawlTasks.slice) > 0 {
		u := sitemapCrawlTasks.slice[0]
		sitemapCrawlTasks.slice = sitemapCrawlTasks.slice[1:len(sitemapCrawlTasks.slice)]
		sitemapCrawlTasks.mu.Unlock()
		completeSitemapCrawl(u)
		return true
	}
	sitemapCrawlTasks.mu.Unlock()
	return false
}

func completeNextPageTask() bool {
	pageCrawlTasks.mu.Lock()
	if len(pageCrawlTasks.slice) > 0 {
		u := pageCrawlTasks.slice[0]
		pageCrawlTasks.slice = pageCrawlTasks.slice[1:len(pageCrawlTasks.slice)]
		pageCrawlTasks.mu.Unlock()
		completePageCrawl(u)
		return true
	}
	pageCrawlTasks.mu.Unlock()
	return false
}

func completeRobotsCrawl(u url) {
	time.Sleep(time.Duration(1000+rand.Int()%1000) * time.Millisecond)
	fmt.Printf("Crawled robots.txt '%s'\n", u.String())
}

func completeSitemapCrawl(u url) {
	time.Sleep(time.Duration(1000+rand.Int()%1000) * time.Millisecond)
	fmt.Printf("Crawled sitemap '%s'\n", u.String())
}

func completePageCrawl(u url) {
	time.Sleep(time.Duration(1000+rand.Int()%1000) * time.Millisecond)
	fmt.Printf("Crawled page '%s'\n", u.String())
}
