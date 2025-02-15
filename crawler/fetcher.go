package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type fetcher interface {
	fetch(url) (string, error)
}

type netFetcher struct{}

func (*netFetcher) fetch(u url) (string, error) {
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

func (d *dummyFetcher) fetch(u url) (string, error) {
	return d.response, nil
}

type fileFetcher struct{}

func (f *fileFetcher) fetch(u url) (string, error) {
	path := u.protocolString() + "/" + u.subdomain + "." + u.domain + "." + u.tld + u.pathString()
	content, err := os.ReadFile("testdata/sites/" + path)
	if err != nil {
		return "", fmt.Errorf("Failed to read file: %v", err)
	}
	return string(content), nil
}
