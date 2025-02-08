package main

import (
	//"fmt"
	"testing"
	"os"
)

func getRobotsTxt(t *testing.T) string {
	content, err := os.ReadFile("testdata/robots.txt")
    if err != nil {
        t.Fatalf("Failed to read file: %v", err)
    }
    return string(content)
}

func TestUrlValidator(t *testing.T) {
    robotsTxt := getRobotsTxt(t)

	validator, _ := processRobotsTxt(robotsTxt)

	testUrls := map[string]bool{
		"/allowed-dir": true,
		"/disallowed" : false,
		"/disallowed/subdir": false,
		"/disallowed/nevermind-this-is-allowed": true,
		"/disallowed/nevermind-this-is-allowed/lol": true,
		"/main.php": false,
		"/endline": false,
		"/endline/not-lol": true,
		"/allowed-dir/disallowed-php.php": false,
	}

	for url, want := range testUrls {
		got := validator.validate(url)
		if got != want {
			t.Errorf("'%s' got %t, want %t", url, got, want)
		}
	}

	//t.Log(blocks)
	//fmt.Println(blocks)
}