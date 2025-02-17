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
	// mu      sync.Mutex
	Tasks   Tasks
	Results Results
	Fetcher fetcher.Fetcher
}

func (r *TaskRunner) Run(concurrency int) {
	go r.Results.ListenToChans()

	loopConcurrently(r.crawlNextRobots, concurrency)
	close(r.Results.RobotsChan)

	<-r.Results.RobotsFinished

	// loopConcurrently(r.crawlNextSitemap, concurrency)
	// close(r.Results.SitemapsChan)

	loopConcurrently(r.crawlNextPage, concurrency)
	close(r.Results.PagesChan)

	<-r.Results.PagesFinished
}

// Attempts to crawl the next robots.txt in the taskList, returns false if no more
func (r *TaskRunner) crawlNextRobots() bool {
	u, ok := r.Tasks.Robots.Next()
	if !ok {
		return false
	}
	res, err := r.Fetcher.Fetch(&u)
	if err != nil {
		robotsResult := RobotsResult{
			Url:           &u,
			Success:       false,
			FailureReason: FetchFailed,
		}

		r.Results.RobotsChan <- &robotsResult
		return true
	}

	validator, _ := robots.Parse(res.Body)
	robotsResult := RobotsResult{
		Url:       &u,
		Success:   true,
		Changed:   true,
		Hash:      hash.Hashs(res.Body),
		Validator: &validator,
	}

	r.Results.RobotsChan <- &robotsResult

	return true
}

// Attempts to crawl the next sitemap in the taskList, returns false if no more
// func (r *TaskRunner) crawlNextSitemap() bool {
// 	u := r.Tasks.Sitemaps.Next()
// 	if u == nil {
// 		return false
// 	}
// 	_, err := r.Fetcher.Fetch(u)
// 	if err != nil {
// 		return true
// 	}
// 	return true
// }

// Attempts to crawl the next page in the taskList, returns false if no more
func (r *TaskRunner) crawlNextPage() bool {
	u, ok := r.Tasks.Pages.Next()
	if !ok {
		return false
	}

	robotsUrl := u
	robotsUrl.Path = []string{"robots.txt"}
	robotsUrl.TrailingSlash = false

	// r.mu.Lock()
	// defer r.mu.Unlock()
	robots, ok := r.Results.Robots[robotsUrl.String()]

	if ok {
		valid := robots.Validator.ValidateUrl(&u)
		if !valid {
			pageResult := PageResult{
				Url:           &u,
				Success:       false,
				FailureReason: RobotsDisallowed,
			}

			r.Results.PagesChan <- &pageResult
			return true
		}
	}

	res, err := r.Fetcher.Fetch(&u)
	if err != nil {
		pageResult := PageResult{
			Url:           &u,
			Success:       false,
			FailureReason: FetchFailed,
		}

		r.Results.PagesChan <- &pageResult
		return true
	}

	p := page.Parse(res.Body, u)
	pageResult := PageResult{
		Url:     &u,
		Success: true,
		Changed: true,
		Page:    &p,
	}

	r.Results.PagesChan <- &pageResult

	return true
}

// Runs n goroutines, each calling t until it return false
func loopConcurrently(t func() bool, workers int) {
	var wg sync.WaitGroup

	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				isNext := t()
				if !isNext {
					break
				}
			}
		}()
	}

	wg.Wait()
}
