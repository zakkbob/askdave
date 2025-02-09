package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"slices"
	"sync"
	"time"
)

var crawledPages = safePageSlice{slice: make([]page, 0)}

var discoveredUrls = safeStringSlice{slice: make([]string, 0)}
var crawledUrls = safeStringSlice{slice: make([]string, 0)}
var uncrawledUrls = safeStringSlice{slice: make([]string, 0)}

var start time.Time

type page struct {
	url         string
	pageTitle     string
	ogTitle       string
	ogDescription string
	ogSiteName    string
	links         []string
}

type safeStringSlice struct {
	mutex sync.Mutex
	slice []string
}

type safePageSlice struct {
	mutex sync.Mutex
	slice []page
}

func getAbsoluteUrl(url string, pageUrl string) (absoluteUrl string) { //Converts relative urls to fully qualified
	if len(url) == 0 { //what the frick happened here
		return url
	}

	if url[0] != '/' {
		return url
	}

	siteUrl := pageUrl

	matched, _ := regexp.MatchString("(?i)https://.*?/.*", pageUrl)

	if matched { //if url pageUrl contains a directory
		siteUrlRegex := regexp.MustCompile("(?i)(https://.*?)/")
		matches := siteUrlRegex.FindStringSubmatch(pageUrl)

		if len(matches) < 2 {
			panic("what the freak!!")
		}

		siteUrl = matches[1]
	}
	absoluteUrl = siteUrl + url

	return absoluteUrl
}

func processDiscoveredUrl(url string) {
	discoveredUrls.mutex.Lock()
	if slices.Contains(discoveredUrls.slice, url) {
		discoveredUrls.mutex.Unlock()
		return
	}
	discoveredUrls.slice = append(discoveredUrls.slice, url)
	discoveredUrls.mutex.Unlock()

	uncrawledUrls.mutex.Lock()
	uncrawledUrls.slice = append(uncrawledUrls.slice, url)
	uncrawledUrls.mutex.Unlock()

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

func extractMetaPropertyContent(page string, metaProperty string) (metaContent string) {
	metaElRegexString := fmt.Sprintf("(?s)<meta[^>]*?property=\"og:%s\"[^>]*?>", metaProperty) //Temporary fix, won't work if content contains a '>'

	metaElRegex := regexp.MustCompile(metaElRegexString)
	metaContentRegex := regexp.MustCompile("(?s)content=\"(.*?)\"")

	elMatches := metaElRegex.FindStringSubmatch(page)
	if len(elMatches) < 1 {
		return ""
	}

	metaEl := elMatches[0]

	contentMatches := metaContentRegex.FindStringSubmatch(metaEl)
	if len(contentMatches) < 2 {
		return ""
	}

	metaContent = contentMatches[1]
	return metaContent
}

func extractPageTitle(page string) (pageTitle string) {
	pageTitleRegex := regexp.MustCompile("(?s)<title.*?>(.*?)</title>") //Temporary, won't match if space are in the tags :(
	matches := pageTitleRegex.FindStringSubmatch(page)

	if len(matches) < 2 {
		return ""
	}

	pageTitle = matches[1]
	return pageTitle
}

func extractPageLinks(page, pageUrl string) (pageLinks []string) {
	linkElRegex := regexp.MustCompile("(?s)<a.*?>") //Wont match if '>' is in the tag somewher :shruggie:
	linkHrefRegex := regexp.MustCompile("(?s)href=\"(.*?)\"")

	elMatches := linkElRegex.FindAllString(page, -1)
	if len(elMatches) < 1 {
		return []string{}
	}

	urlValidationRegex := regexp.MustCompile("https?:\\/\\/(www\\.)?[-a-zA-Z0-9@:%._\\+~#=]{1,256}\\.[a-zA-Z0-9()]{1,6}\\b([-a-zA-Z0-9()@:%_\\+.~#?&//=]*)")

	for _, linkEl := range elMatches {
		hrefMatches := linkHrefRegex.FindStringSubmatch(linkEl)
		if len(hrefMatches) < 1 {
			continue
		}

		linkUrl := getAbsoluteUrl(hrefMatches[1], pageUrl)

		isValidLink := urlValidationRegex.MatchString(linkUrl)
		if !isValidLink {
			//fmt.Printf("Oh darn, not a valid url (%s)\n",linkUrl)
			continue
		}

		pageLinks = append(pageLinks, linkUrl)
	}

	return pageLinks
}

func extractPageData(page, pageUrl string) (pageTitle, ogTitle, ogDescription, ogSiteName string, pageLinks []string) {
	ogTitle = extractMetaPropertyContent(page, "title")
	ogDescription = extractMetaPropertyContent(page, "description")
	ogSiteName = extractMetaPropertyContent(page, "site_name")

	pageTitle = extractPageTitle(page)

	pageLinks = extractPageLinks(page, pageUrl)

	return
}

func fetchPageData(url string) (pageData page) {
	body, _ := fetchPage(url)
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Printf("%s", body)
	// }
	pageTitle, ogTitle, ogDescription, ogSiteName, pageLinks := extractPageData(body, url)
	pageData = page{
		url:           url,
		pageTitle:     pageTitle,
		ogTitle:       ogTitle,
		ogDescription: ogDescription,
		ogSiteName:    ogSiteName,
		links:         pageLinks,
	}

	return pageData
}

func crawlNextUrl() {
	uncrawledUrls.mutex.Lock()
	uncrawledCount := len(uncrawledUrls.slice)
	if len(uncrawledUrls.slice) == 0 {
		uncrawledUrls.mutex.Unlock()
		return
	}

	nextUrl := uncrawledUrls.slice[0]
	uncrawledUrls.slice = slices.Delete(uncrawledUrls.slice, 0, 1)
	uncrawledUrls.mutex.Unlock()

	pageData := fetchPageData(nextUrl)

	for _, url := range pageData.links {
		processDiscoveredUrl(url)
	}

	crawledPages.mutex.Lock()
	crawledPages.slice = append(crawledPages.slice, pageData)
	crawledPages.mutex.Unlock()

	crawledUrls.mutex.Lock()
	crawledCount := len(crawledUrls.slice)
	fmt.Printf("\r%f crawls/s %d Crawled, Uncrawled %d   ", float64(crawledCount)/time.Since(start).Seconds(), crawledCount, uncrawledCount)
	crawledUrls.slice = append(crawledUrls.slice, nextUrl)
	crawledUrls.mutex.Unlock()

}

func autoCrawl(count int) {
	if count < 1 {
		return
	}
	crawlNextUrl()
	autoCrawl(count - 1)
}

func logCrawlStats() {
	crawledUrls.mutex.Lock()
	defer crawledUrls.mutex.Unlock()
	//uncrawledUrls.mutex.Lock()
	//defer uncrawledUrls.mutex.Unlock()

	fmt.Println(len(crawledUrls.slice), "Crawled")
	//fmt.Println(len(uncrawledUrls.slice), "Uncrawled")
}
