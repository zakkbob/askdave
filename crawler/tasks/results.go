package tasks

import (
	"ZakkBob/AskDave/crawler/page"
	"ZakkBob/AskDave/crawler/robots"
)

type Results struct {
	Robots   map[string]robots.UrlValidator
	Sitemaps map[string]string
	Pages    map[string]page.Page
}
