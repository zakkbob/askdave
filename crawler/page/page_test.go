package page_test

import (
	"ZakkBob/AskDave/crawler/fetcher"
	"ZakkBob/AskDave/crawler/page"
	"ZakkBob/AskDave/crawler/url"
	"testing"
)

func TestAddLink(t *testing.T) {
	var p page.Page
	expected := "https://google.com"
	u, _ := url.ParseAbsoluteUrl(expected)
	p.AddLink(u)
	got := p.Links[0].String()
	if got != expected {
		t.Errorf("got '%s', expected '%s'", got, expected)
	}
}

func TestParseBody(t *testing.T) {
	fetcher := &fetcher.FileFetcher{}
	u, _ := url.ParseAbsoluteUrl("https://pagetest.com/index.html")
	p, _ := page.CrawlUrl(u, fetcher)

	link1, _ := url.ParseAbsoluteUrl("https://pagetest.com/example.com")
	link2, _ := url.ParseAbsoluteUrl("https://pagetest.com/lol")

	var hash [16]byte
	copy(hash[:], "2D4659D54B58FD8AB0367C5734AF9632")

	expectedPage := page.Page{
		Url:           u,
		Title:         "Example Page",
		OgTitle:       "og title",
		OgDescription: "og description",
		OgSiteName:    "og sitename",
		Links:         []url.Url{link1, link2},
		Hash:          hash,
	}

	t.Log(expectedPage)
	t.Log(p)
}
