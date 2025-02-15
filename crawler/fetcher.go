package main

import (
	"ZakkBob/AskDave/crawler/urls"
	"fmt"
	"io"
	"net/http"
	"os"
)

type fetcher interface {
	fetch(urls.Url) (string, error)
}

type netFetcher struct{}

func (*netFetcher) fetch(u urls.Url) (string, error) {
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

type dummyFetcher struct {
	response string
}

func (d *dummyFetcher) fetch(u urls.Url) (string, error) {
	return d.response, nil
}

type fileFetcher struct{}

func (f *fileFetcher) fetch(u urls.Url) (string, error) {
	path := u.ProtocolString() + "/" + u.Subdomain() + "." + u.Domain() + "." + u.Tld() + u.PathString()
	content, err := os.ReadFile("testdata/sites/" + path)
	if err != nil {
		return "", fmt.Errorf("Failed to read file: %v", err)
	}
	return string(content), nil
}
