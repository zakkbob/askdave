package main

import (
	"fmt"
	"os"
	"testing"
)

func getRobotsTxt(t *testing.T) string {
	content, err := os.ReadFile("testdata/robots.txt")
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	fmt.Println(content)
	fmt.Println(err)
	return string(content)
}

func TestUrlValidator(t *testing.T) {
	robotsTxt := getRobotsTxt(t)

	validator, _ := ProcessRobotsTxt(robotsTxt)

	testUrls := map[string]bool{
		"/allowed-dir":                              true,
		"/disallowed":                               false,
		"/disallowed/subdir":                        false,
		"/disallowed/nevermind-this-is-allowed":     true,
		"/disallowed/nevermind-this-is-allowed/lol": true,
		"/main.php":                                 false,
		"/endline":                                  false,
		"/endline/not-lol":                          true,
		"/allowed-dir/disallowed-php.php":           false,
		"/allow-subdir/":                            false,
		"/allow-subdir/e":                           true,
	}

	for url, want := range testUrls {
		got := validator.validate(url)
		if got != want {
			t.Errorf("'%s' got %t, want %t", url, got, want)
		}
	}
}
