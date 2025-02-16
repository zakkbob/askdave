//-----------------------------------//
// Processes a webpage into a struct //
//-----------------------------------//

package page

import (
	"ZakkBob/AskDave/crawler/fetcher"
	"ZakkBob/AskDave/crawler/url"
	"crypto/md5"
	"fmt"
	"regexp"
)

type Page struct {
	Url           url.Url
	Title         string
	OgTitle       string
	OgDescription string
	OgSiteName    string
	Links         []url.Url
	Hash          [16]byte //Hash of page (for checking if it has changed)
}

func (p *Page) AddLink(u url.Url) {
	p.Links = append(p.Links, u)
}

func CrawlUrl(url url.Url, fetcher fetcher.Fetcher) (Page, error) {
	body, err := fetcher.Fetch(url)

	if err != nil {
		return Page{}, err
	}

	var p = Page{
		Url:           url,
		Title:         extractTitle(body),
		OgTitle:       extractMeta(body, "title"),
		OgDescription: extractMeta(body, "description"),
		OgSiteName:    extractMeta(body, "site_name"),
		Links:         extractLinks(body, url),
		Hash:          md5.Sum([]byte(body)),
	}
	return p, err
}

func extractMeta(Page string, metaProperty string) string {
	metaElRegexString := fmt.Sprintf("(?s)<meta[^>]*?property=\"og:%s\"[^>]*?>", metaProperty) //Temporary fix, won't work if content contains a '>'

	metaElRegex := regexp.MustCompile(metaElRegexString)
	metaContentRegex := regexp.MustCompile("(?s)content=\"(.*?)\"")

	elMatches := metaElRegex.FindStringSubmatch(Page)
	if len(elMatches) < 1 {
		return ""
	}

	metaEl := elMatches[0]

	contentMatches := metaContentRegex.FindStringSubmatch(metaEl)
	if len(contentMatches) < 2 {
		return ""
	}

	return contentMatches[1]
}

func extractTitle(page string) (pageTitle string) {
	pageTitleRegex := regexp.MustCompile("(?s)<title.*?>(.*?)</title>") //Temporary, won't match if space are in the tags :(
	matches := pageTitleRegex.FindStringSubmatch(page)

	if len(matches) < 2 {
		return ""
	}

	pageTitle = matches[1]
	return pageTitle
}

func extractLinks(body string, pageUrl url.Url) []url.Url {
	var pageLinks []url.Url

	linkElRegex := regexp.MustCompile("(?s)<a.*?>") //Wont match if '>' is in the tag somewher :shruggie:
	linkHrefRegex := regexp.MustCompile("(?s)href=\"(.*?)\"")

	elMatches := linkElRegex.FindAllString(body, -1)
	if len(elMatches) < 1 {
		return []url.Url{}
	}

	for _, linkEl := range elMatches {
		hrefMatches := linkHrefRegex.FindStringSubmatch(linkEl)
		if len(hrefMatches) < 1 {
			continue
		}

		linkUrl, err := url.ParseRelativeUrl(hrefMatches[1], pageUrl)

		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		pageLinks = append(pageLinks, linkUrl)
	}

	return pageLinks
}
