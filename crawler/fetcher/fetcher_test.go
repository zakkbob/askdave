package fetcher_test

import (
	"ZakkBob/AskDave/crawler/fetcher"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFileFetcher(t *testing.T) {
	var f = fetcher.FileFetcher{
		Debug: true,
	}
	got, _ := f.Fetch("https://fetchertest.com/index.html")
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
	got, _ := f.Fetch("fetchertest.com/index.html")
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
