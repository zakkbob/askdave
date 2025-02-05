package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
)

type page struct {
	url            string
	page_title     string
	og_title       string
	og_description string
	og_site_name   string
	links          []string
}

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

func extract_meta_property_content(page string, meta_property string) (meta_content string) {
	meta_el_regex_string := fmt.Sprintf("(?s)<meta[^>]*?property=\"og:%s\"[^>]*?>", meta_property) //Temporary fix, won't work if content contains a '>'

	meta_el_regex := regexp.MustCompile(meta_el_regex_string)
	meta_content_regex := regexp.MustCompile("(?s)content=\"(.*?)\"")

	el_matches := meta_el_regex.FindStringSubmatch(page)
	if len(el_matches) < 1 {
		return ""
	}

	meta_el := el_matches[0]

	content_matches := meta_content_regex.FindStringSubmatch(meta_el)
	if len(content_matches) < 2 {
		return ""
	}

	meta_content = content_matches[1]
	return meta_content
}

func extratct_page_title(page string) (page_title string) {
	page_title_regex := regexp.MustCompile("(?s)<title>(.*?)</title>")
	matches := page_title_regex.FindStringSubmatch(page)

	if len(matches) < 2 {
		return ""
	}

	page_title = matches[1]
	return page_title
}

func extract_page_data(page string) (page_title, og_title, og_description, og_site_name string) {
	og_title = extract_meta_property_content(page, "title")
	og_description = extract_meta_property_content(page, "description")
	og_site_name = extract_meta_property_content(page, "site_name")

	page_title = extratct_page_title(page)

	return
}

func fetch_page_data(url string) (page_data page) {
	body, _ := fetch_page(url)
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Printf("%s", body)
	// }
	page_title, og_title, og_description, og_site_name := extract_page_data(body)
	page_data = page{
		url:            url,
		page_title:     page_title,
		og_title:       og_title,
		og_description: og_description,
		og_site_name:   og_site_name,
	}
	return page_data
}

func main() {
	page_data := fetch_page_data("https://greendungarees.org.uk")
	fmt.Println(page_data)

	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Printf("%s", body)
	// }

	// fmt.Println(page_title)
	// fmt.Println(og_title)
	// fmt.Println(og_description)
	// fmt.Println(og_site_name)
}
