package main

import (
	"ZakkBob/AskDave/crawler/fetcher"
	"ZakkBob/AskDave/crawler/tasks"
	"ZakkBob/AskDave/crawler/url"
	"encoding/json"
	"fmt"
)

func main() {
	r := tasks.TaskRunner{
		Fetcher: &fetcher.FileFetcher{},
		Results: tasks.Results{
			Robots:     make(map[string]*tasks.RobotsResult),
			Pages:      make(map[string]*tasks.PageResult),
			RobotsChan: make(chan *tasks.RobotsResult, 5),
			// SitemapsChan: make(chan *string, 5),
			PagesChan:      make(chan *tasks.PageResult, 5),
			RobotsFinished: make(chan bool, 1),
			PagesFinished:  make(chan bool, 1),
		},
	}

	// for i := 0; i < 2; i++ {
	// 	s := "https://example.com/robots.txt" + strconv.Itoa(i)

	// 	u, _ := url.ParseAbs(s)
	// 	r.Tasks.Robots.Slice = append(r.Tasks.Robots.Slice, u)
	// }

	// for i := 0; i < 2; i++ {
	// 	u, _ := url.ParseAbs(fmt.Sprintf("https://example.com/sitemap%d.xml", i))
	// 	r.Tasks.Sitemaps.Slice = append(r.Tasks.Sitemaps.Slice, u)
	// }

	// for i := 0; i < 2; i++ {
	// 	u, _ := url.ParseAbs(fmt.Sprintf("https://example.com/page%d.html", i))
	// 	r.Tasks.Pages.Slice = append(r.Tasks.Pages.Slice, u)
	// }

	u1, _ := url.ParseAbs("https://taskrunnertest.com/robots.txt")
	r.Tasks.Robots.Slice = append(r.Tasks.Robots.Slice, u1)

	u1, _ = url.ParseAbs("https://taskrunnertest.com/index.html")
	u2, _ := url.ParseAbs("https://taskrunnertest.com/cats.html")
	r.Tasks.Pages.Slice = append(r.Tasks.Pages.Slice, u1, u2)

	r.Run(1)

	j, _ := json.MarshalIndent(r.Results, "", "  ")

	fmt.Println(string(j))
}
