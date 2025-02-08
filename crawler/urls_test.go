package main

import (
	"testing"
)

func TestUrlString(t *testing.T) {
	testUrls := []url{
		url{
			"https", "www", "example", "com",
			[]string{
				"privacy",
			}, true,
		},
		url{
			domain: "example",
			tld:    "com",
		},
	}
	expectedUrls := []string{
		"https://www.example.com/privacy/",
		"example.com",
	}

	for i, u := range testUrls {
		got := u.String()
		expected := expectedUrls[i]
		if got != expected {
			t.Errorf("got '%s', expected '%s'", got, expected)
		}
	}
}
