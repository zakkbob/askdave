package main

import (
	"fmt"
	"strconv"
)

func main() {
	for i := 0; i < 25; i++ {
		s := "https://example.com/robots.txt" + strconv.Itoa(i)

		u, _ := parseAbsoluteUrl(s)
		robotsCrawlTasks.slice = append(robotsCrawlTasks.slice, u)
	}

	for i := 0; i < 25; i++ {
		u, _ := parseAbsoluteUrl(fmt.Sprintf("https://example.com/sitemap%d.xml", i))
		sitemapCrawlTasks.slice = append(sitemapCrawlTasks.slice, u)
	}

	for i := 0; i < 25; i++ {
		u, _ := parseAbsoluteUrl(fmt.Sprintf("https://example.com/page%d.html", i))
		pageCrawlTasks.slice = append(pageCrawlTasks.slice, u)
	}

	loopConcurrently(completeNextRobotsTask, 10)
	loopConcurrently(completeNextSitemapTask, 10)
	loopConcurrently(completeNextPageTask, 10)

	// 	start = time.Now()

	// 	var wg sync.WaitGroup

	// 	discoveredUrls.slice = append(discoveredUrls.slice, "https://mateishome.page")
	// 	uncrawledUrls.slice = append(uncrawledUrls.slice, "https://mateishome.page")

	// 	crawlNextUrl()

	// 	for _ = range 10 {
	// 		wg.Add(1)
	// 		go func() {
	// 			defer wg.Done()
	// 			crawlNextUrl()
	// 		}()
	// 	}

	// 	wg.Wait()

	// 	for _ = range 10 {
	// 		wg.Add(1)
	// 		go func() {
	// 			defer wg.Done()
	// 			autoCrawl(10)
	// 		}()
	// 	}

	// 	wg.Wait()

	// 	//fmt.Println(len(uncrawledUrls.slice), "Uncrawled")

	// 	start = time.Now()
	// 	for _ = range 1000 {
	// 		wg.Add(1)
	// 		go func() {
	// 			defer wg.Done()
	// 			autoCrawl(10000)
	// 		}()
	// 	}

	// 	//go func() {
	// 	//	logCrawlStats()
	// 	//	time.Sleep(100 * time.Millisecond)
	// 	//}()

	// 	wg.Wait()

	// 	fmt.Println(len(crawledUrls.slice), "Crawled")
	// 	fmt.Println(len(uncrawledUrls.slice), "Uncrawled")

	// 	for _, pageData := range crawledPages.slice {
	// 		fmt.Printf("url           - '%s'\n", pageData.url)
	// 		fmt.Printf("pageTitle     - '%s'\n", pageData.pageTitle)
	// 		fmt.Printf("ogTitle       - '%s'\n", pageData.ogTitle)
	// 		fmt.Printf("ogDescription - '%s'\n", pageData.ogDescription)
	// 		fmt.Printf("ogSiteName    - '%s'\n\n", pageData.ogSiteName)
	// 	}

	// //fmt.Println(crawledPages)
	// fmt.Println(len(crawledUrls.slice), "Crawled")
	// fmt.Println(len(uncrawledUrls.slice), "Uncrawled")
}
