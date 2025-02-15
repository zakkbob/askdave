//------------------------------------------------//
// Respresents an interface for fetching webpages //
//------------------------------------------------//

package fetcher

import (
	"ZakkBob/AskDave/crawler/url"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Fetcher interface {
	Fetch(url.Url) (string, error)
}

type NetFetcher struct{}

func (f *NetFetcher) Fetch(u url.Url) (string, error) {
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

type DummyFetcher struct {
	Response string
	Debug    bool
}

// returns DummyFetcher.response, the url parameter is ignored
func (f *DummyFetcher) Fetch(u url.Url) (string, error) {
	if f.Debug {
		fmt.Printf("fetched dummy url '%s'\n", u.String())
	}
	return f.Response, nil
}

type FileFetcher struct {
	Debug bool
}

func (f *FileFetcher) Fetch(u url.Url) (string, error) {
	path := "../testdata/sites/"
	if u.ProtocolString() != "" {
		path += u.ProtocolString() + "/"
	}
	path += u.FQDN() + u.PathString()
	content, err := os.ReadFile(path)
	if f.Debug {
		fmt.Printf("fetching file url '%s'\n", u.String())
		fmt.Printf("path: '%s'\n", path)
		fmt.Printf("content: '%s'\n", content)
	}
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}
	return string(content), nil
}
