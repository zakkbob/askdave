package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"slices"
	"time"
	"sync"
)

type page struct {
	url            string
	page_title     string
	og_title       string
	og_description string
	og_site_name   string
	links          []string
}

type SafeStringSlice struct {
	mutex sync.Mutex
	slice  []string
}

type SafePageSlice struct {
	mutex sync.Mutex
	slice  []page
}

var start time.Time

var crawled_pages = SafePageSlice{slice: make([]page, 0)}

var discovered_urls = SafeStringSlice{slice: make([]string, 0)}
var crawled_urls = SafeStringSlice{slice: make([]string, 0)}
var uncrawled_urls = SafeStringSlice{slice: make([]string, 0)}

func getAbsoluteUrl(url string, page_url string) (absolute_url string) { //Converts relative urls to fully qualified
	if len(url) == 0{ //what the frick happened here
		return url
	}
	
	if url[0] != '/'{
		return url
	}

	site_url := page_url

	matched, _ := regexp.MatchString("(?i)https://.*?/.*", page_url)

	if matched{ //if url page_url contains a directory
		site_url_regex := regexp.MustCompile("(?i)(https://.*?)/")
		matches := site_url_regex.FindStringSubmatch(page_url)

		if len(matches) < 2 {
			panic("what the freak!!")
		}

		site_url = matches[1]
	}
	absolute_url = site_url + url

	return absolute_url
}

func processDiscoveredUrl(url string) {
	discovered_urls.mutex.Lock()
	if slices.Contains(discovered_urls.slice, url){
		discovered_urls.mutex.Unlock()
		return
	}
	discovered_urls.slice = append(discovered_urls.slice, url)
	discovered_urls.mutex.Unlock()

	uncrawled_urls.mutex.Lock()
	uncrawled_urls.slice = append(uncrawled_urls.slice, url)
	uncrawled_urls.mutex.Unlock()

}

func fetchPage(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func extractMetaPropertyContent(page string, meta_property string) (meta_content string) {
	meta_el_regex_string := fmt.Sprintf("(?s)<meta[^>]*?property=\"og:%s\"[^>]*?>", meta_property) //Temporary fix, won't work if content contains a '>'

	meta_el_regex := regexp.MustCompile(meta_el_regex_string)
	meta_content_regex := regexp.MustCompile("(?s)content=\"(.*?)\"")

	el_matches := meta_el_regex.FindStringSubmatch(page)
	if len(el_matches) < 1 {
		return ""
	}

	meta_el := el_matches[0]

	content_matches := meta_content_regex.FindStringSubmatch(meta_el)
	if len(content_matches) < 2 {
		return ""
	}

	meta_content = content_matches[1]
	return meta_content
}

func extractPageTitle(page string) (page_title string) {
	page_title_regex := regexp.MustCompile("(?s)<title.*?>(.*?)</title>") //Temporary, won't match if space are in the tags :(
	matches := page_title_regex.FindStringSubmatch(page)

	if len(matches) < 2 {
		return ""
	}

	page_title = matches[1]
	return page_title
}

func extractPageLinks(page, page_url string) (page_links []string) {
	link_el_regex   := regexp.MustCompile("(?s)<a.*?>") //Wont match if '>' is in the tag somewher :shruggie:
	link_href_regex := regexp.MustCompile("(?s)href=\"(.*?)\"") 

	el_matches := link_el_regex.FindAllString(page, -1)
	if len(el_matches) < 1 {
		return []string{}
	}

	url_validation_regex := regexp.MustCompile("https?:\\/\\/(www\\.)?[-a-zA-Z0-9@:%._\\+~#=]{1,256}\\.[a-zA-Z0-9()]{1,6}\\b([-a-zA-Z0-9()@:%_\\+.~#?&//=]*)")

	for _, link_el := range el_matches {
		href_matches := link_href_regex.FindStringSubmatch(link_el)
		if len(href_matches) < 1 {
			continue
		}
		
		link_url := getAbsoluteUrl(href_matches[1], page_url)

		is_valid_link := url_validation_regex.MatchString(link_url)
		if (!is_valid_link) {
			//fmt.Printf("Oh darn, not a valid url (%s)\n",link_url)
			continue
		}

		page_links = append(page_links, link_url)
	}

	return page_links
}

func extractPageData(page, page_url string) (page_title, og_title, og_description, og_site_name string, page_links []string) {
	og_title = extractMetaPropertyContent(page, "title")
	og_description = extractMetaPropertyContent(page, "description")
	og_site_name = extractMetaPropertyContent(page, "site_name")

	page_title = extractPageTitle(page)

	page_links = extractPageLinks(page, page_url)

	return
}

func fetchPageData(url string) (page_data page) {
	body, _ := fetchPage(url)
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Printf("%s", body)
	// }
	page_title, og_title, og_description, og_site_name, page_links := extractPageData(body, url)
	page_data = page{
		url:            url,
		page_title:     page_title,
		og_title:       og_title,
		og_description: og_description,
		og_site_name:   og_site_name,
		links:          page_links,
	}

	return page_data
}

func crawlNextUrl() {
	uncrawled_urls.mutex.Lock()
	uncrawled_count := len(uncrawled_urls.slice)
	if len(uncrawled_urls.slice) == 0{
		uncrawled_urls.mutex.Unlock()
		return
	}

	next_url := uncrawled_urls.slice[0]
	uncrawled_urls.slice = slices.Delete(uncrawled_urls.slice, 0, 1)
	uncrawled_urls.mutex.Unlock()

	page_data := fetchPageData(next_url)

	for _, url := range page_data.links {
		processDiscoveredUrl(url)
	}

	crawled_pages.mutex.Lock()
	crawled_pages.slice = append(crawled_pages.slice, page_data)
	crawled_pages.mutex.Unlock()

	crawled_urls.mutex.Lock()
	crawled_count := len(crawled_urls.slice)
	fmt.Printf("\r%f crawls/s %d Crawled, Uncrawled %d   ", float64(crawled_count)/time.Since(start).Seconds(), crawled_count, uncrawled_count)
	crawled_urls.slice = append(crawled_urls.slice, next_url)
	crawled_urls.mutex.Unlock()

}

func autoCrawl(count int) {
	if count < 1 {
		return
	}
	crawlNextUrl()
	autoCrawl(count-1)
}

func logCrawlStats() {
	crawled_urls.mutex.Lock()
	defer crawled_urls.mutex.Unlock()
	//uncrawled_urls.mutex.Lock()
	//defer uncrawled_urls.mutex.Unlock()

	fmt.Println(len(crawled_urls.slice), "Crawled")
	//fmt.Println(len(uncrawled_urls.slice), "Uncrawled")
}

func main() {
	start = time.Now()

	var wg sync.WaitGroup

	discovered_urls.slice = append(discovered_urls.slice, "https://mateishome.page")
	uncrawled_urls.slice = append(uncrawled_urls.slice, "https://mateishome.page")

	crawlNextUrl()

	for _ = range 10{
		wg.Add(1)
		go func() {
			defer wg.Done()
			crawlNextUrl()
		}()
	}

	wg.Wait()

	for _ = range 10{
		wg.Add(1)
		go func() {
			defer wg.Done()
			autoCrawl(10)
		}()
	}

	wg.Wait()

	//fmt.Println(len(uncrawled_urls.slice), "Uncrawled")

	start = time.Now()
	for _ = range 1000{
		wg.Add(1)
		go func() {
			defer wg.Done()
			autoCrawl(10000)
		}()
	}

	//go func() {
	//	logCrawlStats()
	//	time.Sleep(100 * time.Millisecond)
	//}()

	wg.Wait()

	fmt.Println(len(crawled_urls.slice), "Crawled")
	fmt.Println(len(uncrawled_urls.slice), "Uncrawled")



	for _, page_data := range crawled_pages.slice{
		fmt.Printf("url            - '%s'\n", page_data.url)
		fmt.Printf("page_title     - '%s'\n", page_data.page_title)
		fmt.Printf("og_title       - '%s'\n", page_data.og_title)
		fmt.Printf("og_description - '%s'\n", page_data.og_description)
		fmt.Printf("og_site_name   - '%s'\n\n", page_data.og_site_name)
	}

	//fmt.Println(crawled_pages)
	fmt.Println(len(crawled_urls.slice), "Crawled")
	fmt.Println(len(uncrawled_urls.slice), "Uncrawled")
}
