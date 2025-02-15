package main

import (
	"ZakkBob/AskDave/crawler/tasks"
	"ZakkBob/AskDave/crawler/url"
	"fmt"
	"strconv"
)

func main() {
	var r tasks.TaskRunner

	for i := 0; i < 25; i++ {
		s := "https://example.com/robots.txt" + strconv.Itoa(i)

		u, _ := url.ParseAbsoluteUrl(s)
		r.Tasks.Robots = append(r.Tasks.Robots, u)
	}

	for i := 0; i < 25; i++ {
		u, _ := url.ParseAbsoluteUrl(fmt.Sprintf("https://example.com/sitemap%d.xml", i))
		r.Tasks.Sitemaps = append(r.Tasks.Sitemaps, u)
	}

	for i := 0; i < 25; i++ {
		u, _ := url.ParseAbsoluteUrl(fmt.Sprintf("https://example.com/page%d.html", i))
		r.Tasks.Pages = append(r.Tasks.Pages, u)
	}

	r.Run(25)
}
