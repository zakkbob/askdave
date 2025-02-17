//-------------------------------------------------------------//
// An interface for fetching webpages and some implementations //
//-------------------------------------------------------------//

package fetcher

import (
	"ZakkBob/AskDave/crawler/url"
	"embed"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

//go:embed testsites/*
var files embed.FS

type Fetcher interface {
	Fetch(*url.Url) (Response, error)
}

type Response struct {
	Body       string
	StatusCode int
}

type NetFetcher struct{}

func (f *NetFetcher) Fetch(u *url.Url) (Response, error) {
	resp, err := http.Get(u.String())
	if err != nil {
		return Response{}, err // Get was unsuccessful, url probably doesnt exist or something, who knows
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}

	return Response{
		Body:       string(body),
		StatusCode: resp.StatusCode,
	}, nil
}

type DummyFetcher struct {
	Response  string
	Delay     time.Duration
	RandDelay time.Duration
	Debug     bool
}

func (f *DummyFetcher) sleep() time.Duration {
	if f.Delay == 0 && f.RandDelay == 0 {
		return 0
	}

	if f.RandDelay == 0 {
		time.Sleep(f.Delay)
		return f.Delay
	}

	extraDelay := time.Duration(rand.Int() % (int(f.RandDelay)))
	totalDelay := f.Delay + extraDelay
	time.Sleep(totalDelay)
	return totalDelay
}

// returns DummyFetcher.response, the url parameter is ignored
func (f *DummyFetcher) Fetch(u *url.Url) (Response, error) {
	totalDelay := f.sleep()

	if f.Debug {
		fmt.Printf("fetched dummy url '%s' after %s\n", u.String(), totalDelay.String())
	}
	return Response{
		Body:       f.Response,
		StatusCode: 200,
	}, nil
}

type FileFetcher struct {
	Debug bool
}

func (f *FileFetcher) Fetch(u *url.Url) (Response, error) {
	path := "testsites/"
	path += u.FQDN() + u.PathString()
	content, err := files.ReadFile(path)
	if f.Debug {
		fmt.Printf("fetching file url '%s'\n", u.String())
		fmt.Printf("path: '%s'\n", path)
		fmt.Printf("content: '%s'\n", content)
	}
	if err != nil {
		return Response{}, fmt.Errorf("failed to read file: %v", err)
	}
	return Response{
		Body:       string(content),
		StatusCode: 200,
	}, nil
}
