package tasks

import (
	"ZakkBob/AskDave/crawler/hash"
	"ZakkBob/AskDave/crawler/page"
	"ZakkBob/AskDave/crawler/robots"
	"ZakkBob/AskDave/crawler/url"
)

type FailureReason = int

const (
	NoFailure FailureReason = iota
	RobotsDisallowed
)

type Results struct {
	Robots   map[string]RobotsResult
	Sitemaps map[string]string
	Pages    map[string]PageResult
}

type PageResult struct {
	Url           url.Url
	Success       bool
	FailureReason FailureReason
	Changed       bool
	Page          page.Page
}

type RobotsResult struct {
	Url       url.Url
	Hash      hash.Hash
	Changed   bool
	Validator robots.UrlValidator
}
