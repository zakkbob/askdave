package main

import (
	"ZakkBob/AskDave/crawler/fetcher"
	"ZakkBob/AskDave/crawler/tasks"
	"encoding/json"
	"fmt"
)

func main() {
	r := tasks.TaskRunner{
		Fetcher: &fetcher.NetFetcher{},
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

	// 	data :=
	// 		`{
	//   "robots": {
	//     "slice": [
	//       "https://taskrunnertest.com/robots.txt/e"
	//     ]
	//   },
	//   "sitemaps": {
	//     "slice": null
	//   },
	//   "pages": {
	//     "slice": [
	//       "https://taskrunnertest.com/index.html",
	//       "https://taskrunnertest.com/cats.html",
	// 	  "https://taskrunnertest.com/notfound.php",
	// 	  "https://taskrunnertest.com/disallowed/secrets.txt"
	//     ]
	//   }
	// }`

	data :=
		`{
"robots": {
"slice": null
},
"sitemaps": {
"slice": null
},
"pages": {
"slice": [
"https://mateishome.page"
]
}
}`

	var t tasks.Tasks
	json.Unmarshal([]byte(data), &t)

	r.Tasks.Pages = t.Pages
	r.Tasks.Robots = t.Robots

	r.Run(5)

	j, _ := json.MarshalIndent(&r.Results, "", "  ")

	fmt.Println(string(j))
}
