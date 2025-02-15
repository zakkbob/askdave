//-------------------------------//
// Works through all crawl tasks //
//-------------------------------//

package tasks

import (
	"ZakkBob/AskDave/crawler/url"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type TaskRunner struct {
	Tasks   Tasks
	Results Results
}

func (r *TaskRunner) Run(concurrency int) {
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

func (runner *TaskRunner) completeNextRobotsTask() bool {
	runner.Tasks.RobotsMu.Lock()
	if len(runner.Tasks.Robots) > 0 {
		u := runner.Tasks.Robots[0]
		runner.Tasks.Robots = runner.Tasks.Robots[1:len(runner.Tasks.Robots)]
		runner.Tasks.RobotsMu.Unlock()
		completeRobotsCrawl(u)
		return true
	}
	runner.Tasks.RobotsMu.Unlock()
	return false
}

func (r *TaskRunner) completeNextSitemapTask() bool {
	r.Tasks.SitemapsMu.Lock()
	if len(r.Tasks.Sitemaps) > 0 {
		u := r.Tasks.Sitemaps[0]
		r.Tasks.Sitemaps = r.Tasks.Sitemaps[1:len(r.Tasks.Sitemaps)]
		r.Tasks.SitemapsMu.Unlock()
		completeSitemapCrawl(u)
		return true
	}
	r.Tasks.SitemapsMu.Unlock()
	return false
}

func (r *TaskRunner) completeNextPageTask() bool {
	r.Tasks.PagesMu.Lock()
	if len(r.Tasks.Pages) > 0 {
		u := r.Tasks.Pages[0]
		r.Tasks.Pages = r.Tasks.Pages[1:len(r.Tasks.Pages)]
		r.Tasks.PagesMu.Unlock()
		completePageCrawl(u)
		return true
	}
	r.Tasks.PagesMu.Unlock()
	return false
}

func completeRobotsCrawl(u url.Url) {
	time.Sleep(time.Duration(1000+rand.Int()%1000) * time.Millisecond)
	fmt.Printf("Crawled robots.txt '%s'\n", u.String())
}

func completeSitemapCrawl(u url.Url) {
	time.Sleep(time.Duration(1000+rand.Int()%1000) * time.Millisecond)
	fmt.Printf("Crawled sitemap '%s'\n", u.String())
}

func completePageCrawl(u url.Url) {
	time.Sleep(time.Duration(1000+rand.Int()%1000) * time.Millisecond)
	fmt.Printf("Crawled page '%s'\n", u.String())
}
