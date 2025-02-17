package fetcher_test

import (
	"ZakkBob/AskDave/crawler/fetcher"
	"ZakkBob/AskDave/crawler/url"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFileFetcher(t *testing.T) {
	var f = fetcher.FileFetcher{
		Debug: true,
	}
	s, _ := url.ParseAbs("https://fetchertest.com/index.html")
	got, _ := f.Fetch(&s)
	want := "Success"
	if got.Body != want {
		t.Errorf("got '%s', wanted '%s'", got.Body, want)
	}
}

func TestDummyFetcher(t *testing.T) {
	var f = fetcher.DummyFetcher{
		Debug:    true,
		Response: "Success",
	}
	s, _ := url.ParseAbs("fetchertest.com/index.html")
	got, _ := f.Fetch(&s)
	want := "Success"
	if got.Body != want {
		t.Errorf("got '%s', wanted '%s'", got.Body, want)
	}
}

func TestNetFetcher(t *testing.T) {
	mockHandler := http.NewServeMux()
	mockHandler.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Success"))
	})

	mockServer := httptest.NewServer(mockHandler)
	defer mockServer.Close()

	fetcher := &fetcher.NetFetcher{}

	want := "Success"
	got, err := fetcher.Fetch(mockServer.URL + "/test")
	if err != nil {
		t.Fatalf("expected no error, but got: %v", err)
	}

	if got.Body != want {
		t.Errorf("got '%s', wanted '%s'", got.Body, want)
	}
}
