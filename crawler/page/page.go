//-----------------------------------//
// Processes a webpage into a struct //
//-----------------------------------//

package page

import (
	"ZakkBob/AskDave/crawler/url"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

type Page struct {
	u             url.Url
	pageTitle     string
	ogTitle       string
	ogDescription string
	ogSiteName    string
	links         []url.Url
}

func (p *Page) addLink(u url.Url) {
	p.links = append(p.links, u)
}

func crawlPage(u url.Url) (Page, error) {
	b, err := fetchPageBody(u)

	if err != nil {
		return Page{u: u}, err
	}

	p, err := parseBody(b, u)
	p.u = u

	return p, err
}

func parseBody(b string, u url.Url) (Page, error) {
	var p Page
	p.ogTitle = extractBodyMeta(b, "title")
	p.ogDescription = extractBodyMeta(b, "description")
	p.ogSiteName = extractBodyMeta(b, "site_name")

	p.pageTitle = extractBodyTitle(b)

	p.links = extractBodyLinks(b, u)

	return p, nil
}

func fetchPageBody(u url.Url) (string, error) {
	resp, err := http.Get(u.String())
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

func extractBodyMeta(Page string, metaProperty string) (metaContent string) {
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

	metaContent = contentMatches[1]
	return metaContent
}

func extractBodyTitle(page string) (pageTitle string) {
	pageTitleRegex := regexp.MustCompile("(?s)<title.*?>(.*?)</title>") //Temporary, won't match if space are in the tags :(
	matches := pageTitleRegex.FindStringSubmatch(page)

	if len(matches) < 2 {
		return ""
	}

	pageTitle = matches[1]
	return pageTitle
}

func extractBodyLinks(body string, pageUrl url.Url) (pageLinks []url.Url) {
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
