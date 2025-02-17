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
		Fetcher: &fetcher.NetFetcher{},
		Results: tasks.Results{
			RobotsChan:   make(chan *tasks.RobotsResult, 5),
			SitemapsChan: make(chan *string, 5),
			PagesChan:    make(chan *tasks.PageResult, 5),
			Finished:     make(chan bool, 1),
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

	u2, _ := url.ParseAbs("https://emateishome.page")
	r.Tasks.Pages.Slice = append(r.Tasks.Pages.Slice, u2)

	r.Run(1)

	j, _ := json.MarshalIndent(r.Results, "", "  ")

	fmt.Println(string(j))
}
