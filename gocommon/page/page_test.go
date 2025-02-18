package page_test

import (
	"testing"

	"github.com/ZakkBob/AskDave/gocommon/page"
	"github.com/ZakkBob/AskDave/gocommon/url"
)

func TestAddLink(t *testing.T) {
	var p page.Page
	expected := "https://google.com"
	u, _ := url.ParseAbs(expected)
	p.AddLink(u)
	got := p.Links[0].String()
	if got != expected {
		t.Errorf("got '%s', expected '%s'", got, expected)
	}
}

func TestParseBody(t *testing.T) {
	b := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Example Page</title>
    <meta property="og:site_name" content="og sitename">
    <meta property="og:description" content="og description">
    <meta property="og:title" content="og title">
</head>
<body>
    <a href="example.com"></a>
    <a href="./lol"></a>
</body>
</html>`
	u, _ := url.ParseAbs("https://pagetest.com/index.html")
	p := page.Parse(b, u)

	link1, _ := url.ParseAbs("https://pagetest.com/example.com")
	link2, _ := url.ParseAbs("https://pagetest.com/lol")

	expectedHash := [16]byte{45, 70, 89, 213, 75, 88, 253, 138, 176, 54, 124, 87, 52, 175, 150, 50}

	expectedPage := page.Page{
		Url:           u,
		Title:         "Example Page",
		OgTitle:       "og title",
		OgDescription: "og description",
		OgSiteName:    "og sitename",
		Links:         []url.Url{link1, link2},
		Hash:          expectedHash,
	}

	if expectedPage.Url.String() != p.Url.String() {
		t.Errorf("got url '%s', expected '%s'", expectedPage.Url.String(), p.Url.String())
	}
	if expectedPage.Title != p.Title {
		t.Errorf("got title '%s', expected '%s'", expectedPage.Title, p.Title)
	}
	if expectedPage.OgTitle != p.OgTitle {
		t.Errorf("got ogTitle '%s', expected '%s'", expectedPage.OgTitle, p.OgTitle)
	}
	if expectedPage.OgDescription != p.OgDescription {
		t.Errorf("got ogDescription '%s', expected '%s'", expectedPage.OgDescription, p.OgDescription)
	}
	if expectedPage.OgSiteName != p.OgSiteName {
		t.Errorf("got ogSiteName '%s', expected '%s'", expectedPage.OgSiteName, p.OgSiteName)
	}
	for i, got := range p.Links {
		want := expectedPage.Links[i]
		if got.String() != want.String() {
			t.Errorf("got links '%+v', expected '%v'", expectedPage.Links, p.Links)
			break
		}
	}
	if expectedPage.Hash != p.Hash {
		t.Errorf("got hash '%s', expected '%s'", expectedPage.Hash.String(), p.Hash.String())
	}
}
