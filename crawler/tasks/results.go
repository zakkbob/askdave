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
	FetchFailed
)

type Results struct {
	Robots map[string]*RobotsResult `json:"robots,omitempty"`
	// Sitemaps     map[string]*string       `json:"sitemaps,omitempty"`
	Pages      map[string]*PageResult `json:"pages,omitempty"`
	RobotsChan chan *RobotsResult     `json:"-"`
	// SitemapsChan chan *string           `json:"-"`
	PagesChan      chan *PageResult `json:"-"`
	RobotsFinished chan bool        `json:"-"`
	PagesFinished  chan bool        `json:"-"`
}

type PageResult struct {
	Url           *url.Url      `json:"url"`
	Success       bool          `json:"success"`
	FailureReason FailureReason `json:"failure_reason,omitempty"`
	Changed       bool          `json:"changed,omitempty"`
	Page          *page.Page    `json:"page,omitempty"`
}

type RobotsResult struct {
	Url           *url.Url             `json:"robots"`
	Success       bool                 `json:"success"`
	FailureReason FailureReason        `json:"failure_reason,omitempty"`
	Hash          hash.Hash            `json:"hash,omitempty"`
	Changed       bool                 `json:"changed,omitempty"`
	Validator     *robots.UrlValidator `json:"validator,omitempty"`
}

func (r *Results) ListenToChans() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			robotsResult, ok := <-r.RobotsChan
			if !ok && robotsResult == nil {
				r.RobotsFinished <- true
				return
			}
			r.Robots[robotsResult.Url.String()] = robotsResult
		}
	}()
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	for {
	// 		sitemapResult, ok := <-r.SitemapsChan
	// 		if !ok && sitemapResult == nil {
	// 			return
	// 		}
	// 		r.Sitemaps[sitemapResult.Url.String()] = sitemapResult
	// 	}
	// }()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			pageResult, ok := <-r.PagesChan
			if !ok && pageResult == nil {
				r.PagesFinished <- true
				return
			}
			r.Pages[pageResult.Url.String()] = pageResult
		}
	}()
	wg.Wait()
}
