package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
)

func fetch_page(url string) (string, error) {
	resp, err := http.Get(url)
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

func extract_page_info(page string) (page_title, og_title, og_description string) {
	title_regex, _ := regexp.Compile("(?s)property=\"og:title\".*?content=\"(.*?)\"")
	matches := title_regex.FindStringSubmatch(page)
	page_title = matches[0]

	og_title = ""
	og_description = ""

	return
}

func main() {
	url := "https://greendungarees.org.uk"
	body, _ := fetch_page(url)
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Printf("%s", body)
	// }
	page_title, _, _ := extract_page_info(body)
	fmt.Println(page_title)
}
