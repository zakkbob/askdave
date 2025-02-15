//-------------------------------//
// Works through all crawl tasks //
//-------------------------------//

package main

import (
	"ZakkBob/AskDave/crawler/robots"
	"ZakkBob/AskDave/crawler/urls"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type tasks struct {
	robots   safeUrlSlice
	sitemaps safeUrlSlice
	pages    safeUrlSlice
}

type results struct {
	robots   map[string]robots.UrlValidator
	sitemaps map[string]string
	pages    map[string]page
}

type taskRunner struct {
	t tasks
	r results
}

func (r *taskRunner) run(concurrency int) {
	loopConcurrently(r.completeNextRobotsTask, concurrency)
	loopConcurrently(r.completeNextSitemapTask, concurrency)
	loopConcurrently(r.completeNextPageTask, concurrency)
}

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

func (r *taskRunner) completeNextRobotsTask() bool {
	r.t.robots.mu.Lock()
	if len(r.t.robots.slice) > 0 {
		u := r.t.robots.slice[0]
		r.t.robots.slice = r.t.robots.slice[1:len(r.t.robots.slice)]
		r.t.robots.mu.Unlock()
		completeRobotsCrawl(u)
		return true
	}
	r.t.robots.mu.Unlock()
	return false
}

func (r *taskRunner) completeNextSitemapTask() bool {
	r.t.sitemaps.mu.Lock()
	if len(r.t.sitemaps.slice) > 0 {
		u := r.t.sitemaps.slice[0]
		r.t.sitemaps.slice = r.t.sitemaps.slice[1:len(r.t.sitemaps.slice)]
		r.t.sitemaps.mu.Unlock()
		completeSitemapCrawl(u)
		return true
	}
	r.t.sitemaps.mu.Unlock()
	return false
}

func (r *taskRunner) completeNextPageTask() bool {
	r.t.pages.mu.Lock()
	if len(r.t.pages.slice) > 0 {
		u := r.t.pages.slice[0]
		r.t.pages.slice = r.t.pages.slice[1:len(r.t.pages.slice)]
		r.t.pages.mu.Unlock()
		completePageCrawl(u)
		return true
	}
	r.t.pages.mu.Unlock()
	return false
}

func completeRobotsCrawl(u urls.Url) {
	time.Sleep(time.Duration(1000+rand.Int()%1000) * time.Millisecond)
	fmt.Printf("Crawled robots.txt '%s'\n", u.String())
}

func completeSitemapCrawl(u urls.Url) {
	time.Sleep(time.Duration(1000+rand.Int()%1000) * time.Millisecond)
	fmt.Printf("Crawled sitemap '%s'\n", u.String())
}

func completePageCrawl(u urls.Url) {
	time.Sleep(time.Duration(1000+rand.Int()%1000) * time.Millisecond)
	fmt.Printf("Crawled page '%s'\n", u.String())
}
