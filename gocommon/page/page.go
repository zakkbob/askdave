//-----------------------------------//
// Processes a webpage into a struct //
//-----------------------------------//

package page

import (
	"fmt"
	"regexp"

	"github.com/ZakkBob/AskDave/gocommon/url"

	"github.com/ZakkBob/AskDave/gocommon/hash"
)

type Page struct {
	Url           url.URL   `json:"url"`
	Title         string    `json:"title"`
	OgTitle       string    `json:"og_title"`
	OgDescription string    `json:"og_description"`
	OgSiteName    string    `json:"og_site_name"`
	Links         []url.URL `json:"links"`
	Hash          hash.Hash `json:"hash"`
}

func (p *Page) AddLink(u url.URL) {
	p.Links = append(p.Links, u)
}

// Parses webpage body (string) into a page struct
func Parse(body string, u url.URL) Page {
	var p = Page{
		Url:           u,
		Title:         extractTitle(body),
		OgTitle:       extractMeta(body, "title"),
		OgDescription: extractMeta(body, "description"),
		OgSiteName:    extractMeta(body, "site_name"),
		Links:         extractLinks(body, u),
		Hash:          hash.Hashs(body),
	}
	return p
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

func extractLinks(body string, pageUrl url.URL) []url.URL {
	var pageLinks []url.URL

	linkElRegex := regexp.MustCompile("(?s)<a.*?>") //Wont match if '>' is in the tag somewher :shruggie:
	linkHrefRegex := regexp.MustCompile("(?s)href=\"(.*?)\"")

	elMatches := linkElRegex.FindAllString(body, -1)
	if len(elMatches) < 1 {
		return []url.URL{}
	}

	for _, linkEl := range elMatches {
		hrefMatches := linkHrefRegex.FindStringSubmatch(linkEl)
		if len(hrefMatches) < 1 {
			continue
		}

		relativeUrl, _ := url.Parse(hrefMatches[1])
		linkUrl := pageUrl.ResolveReference(relativeUrl)

		pageLinks = append(pageLinks, *linkUrl)
	}

	return pageLinks
}
