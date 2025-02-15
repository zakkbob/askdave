package tasks

import (
	"ZakkBob/AskDave/crawler/url"
	"sync"
)

type Tasks struct {
	Robots     []url.Url
	RobotsMu   sync.Mutex //Should be able to remove this eventually (bit clunky)
	Sitemaps   []url.Url
	SitemapsMu sync.Mutex
	Pages      []url.Url
	PagesMu    sync.Mutex
}
