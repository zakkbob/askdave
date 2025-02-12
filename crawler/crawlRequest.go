package main

type crawlRequestType int

const (
    pageCrawlRequest crawlRequestType = iota
    robotsCrawlRequest
    sitemapCrawlRequest
)

type crawlRequest struct {
	crawlType crawlRequestType
	url url
}

