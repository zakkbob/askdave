package main

import (
	"testing"
)

func TestUrlString(t *testing.T) {
	testUrls := []url{
		url{
			"https", "www", "example", "com", 8080,
			[]string{
				"privacy",
			}, true,
		},
		url{
			domain: "example",
			tld:    "com",
		},
		url{
			protocol: "http",
			domain:"example",
			tld:  "com",
			port:  80,
		},
	}
	expectedUrls := []string{
		"https://www.example.com:8080/privacy/",
		"example.com",
		"http://example.com",
	}

	for i, u := range testUrls {
		got := u.String()
		expected := expectedUrls[i]
		if got != expected {
			t.Errorf("got '%s', expected '%s'", got, expected)
		}
	}
}
