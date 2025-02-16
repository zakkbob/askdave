//-------------------------------//
// Works through all crawl tasks //
//-------------------------------//

package tasks

import (
	"ZakkBob/AskDave/crawler/fetcher"
	"ZakkBob/AskDave/crawler/hash"
	"ZakkBob/AskDave/crawler/page"
	"ZakkBob/AskDave/crawler/robots"
	"sync"
)

type TaskRunner struct {
	Tasks   Tasks
	Results Results
	Fetcher fetcher.Fetcher
}

func (r *TaskRunner) Run(concurrency int) {
	go r.Results.ListenToChans()

	loopConcurrently(r.crawlNextRobots, concurrency)
	close(r.Results.RobotsChan)

	loopConcurrently(r.crawlNextSitemap, concurrency)
	close(r.Results.SitemapsChan)

	loopConcurrently(r.crawlNextPage, concurrency)
	close(r.Results.PagesChan)

	<-r.Results.Finished
}

// Attempts to crawl the next robots.txt in the taskList, returns false if no more
func (r *TaskRunner) crawlNextRobots() (bool, error) {
	u := r.Tasks.Robots.Next()
	if u == nil {
		return false, nil
	}
	b, err := r.Fetcher.Fetch(u)
	if err != nil {
		return true, err
	}

	validator, _ := robots.Parse(b)
	robotsResult := RobotsResult{
		Url:       u,
		Changed:   true,
		Hash:      hash.Hashs(b),
		Validator: &validator,
	}

	r.Results.RobotsChan <- &robotsResult

	return true, nil
}

// Attempts to crawl the next sitemap in the taskList, returns false if no more
func (r *TaskRunner) crawlNextSitemap() (bool, error) {
	u := r.Tasks.Sitemaps.Next()
	if u == nil {
		return false, nil
	}
	_, err := r.Fetcher.Fetch(u)
	if err != nil {
		return true, err
	}
	return true, nil
}

// Attempts to crawl the next page in the taskList, returns false if no more
func (r *TaskRunner) crawlNextPage() (bool, error) {
	u := r.Tasks.Pages.Next()
	if u == nil {
		return false, nil
	}
	b, err := r.Fetcher.Fetch(u)
	if err != nil {
		return true, err
	}

	p := page.Parse(b, *u)
	pageResult := PageResult{
		Url:     u,
		Success: true,
		Changed: true,
		Page:    &p,
	}

	r.Results.PagesChan <- &pageResult

	return true, nil
}

// Runs n goroutines, each calling t until it return false
func loopConcurrently(t func() (bool, error), workers int) {
	var wg sync.WaitGroup

	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				isNext, err := t()
				if err != nil {
					panic(err.Error())
				}
				if !isNext {
					break
				}
			}
		}()
	}

	wg.Wait()
}
