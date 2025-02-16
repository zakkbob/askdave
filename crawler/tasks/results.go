package tasks

import (
	"ZakkBob/AskDave/crawler/hash"
	"ZakkBob/AskDave/crawler/page"
	"ZakkBob/AskDave/crawler/robots"
	"ZakkBob/AskDave/crawler/url"
	"sync"
)

type FailureReason = int

const (
	NoFailure FailureReason = iota
	RobotsDisallowed
)

type Results struct {
	Robots       []RobotsResult     `json:"robots"`
	Sitemaps     []string           `json:"sitemaps"`
	Pages        []PageResult       `json:"pages"`
	RobotsChan   chan *RobotsResult `json:"-"`
	SitemapsChan chan *string       `json:"-"`
	PagesChan    chan *PageResult   `json:"-"`
	Finished     chan bool          `json:"-"`
}

type PageResult struct {
	Url           *url.Url      `json:"url"`
	Success       bool          `json:"success"`
	FailureReason FailureReason `json:"failure_reason,omitempty"`
	Changed       bool          `json:"changed"`
	Page          *page.Page    `json:"page"`
}

type RobotsResult struct {
	Url       *url.Url             `json:"robots"`
	Hash      hash.Hash            `json:"hash"`
	Changed   bool                 `json:"changed"`
	Validator *robots.UrlValidator `json:"validator"`
}

func (r *Results) ListenToChans() {
	var wg sync.WaitGroup

	wg.Add(3)
	go func() {
		defer wg.Done()
		for {
			robotsResult, ok := <-r.RobotsChan
			if !ok && robotsResult == nil {
				return
			}
			r.Robots = append(r.Robots, *robotsResult)
		}
	}()
	go func() {
		defer wg.Done()
		for {
			sitemapResult, ok := <-r.SitemapsChan
			if !ok && sitemapResult == nil {
				return
			}
			r.Sitemaps = append(r.Sitemaps, *sitemapResult)
		}
	}()
	go func() {
		defer wg.Done()
		for {
			pageResult, ok := <-r.PagesChan
			if !ok && pageResult == nil {
				return
			}
			r.Pages = append(r.Pages, *pageResult)
		}
	}()
	wg.Wait()
	r.Finished <- true
}
