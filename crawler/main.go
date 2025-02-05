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

func extract_meta_property_content(page string, meta_property string) (meta_content string) {
	meta_el_regex_string := fmt.Sprintf("(?s)<meta[^>]*?property=\"og:%s\"[^>]*?>", meta_property) //Temporary fix, won't work if content contains a '>'
	
	meta_el_regex      := regexp.MustCompile(meta_el_regex_string)
	meta_content_regex := regexp.MustCompile("(?s)content=\"(.*?)\"")
	
	el_matches := meta_el_regex.FindStringSubmatch(page)
	if len(el_matches) < 1{
		return ""
	}

	meta_el := el_matches[0]

	content_matches := meta_content_regex.FindStringSubmatch(meta_el)
	if len(content_matches) < 2{
		return ""
	}

	meta_content = content_matches[1]
	return meta_content
}

func extract_page_info(page string) (page_title, og_title, og_description, og_site_name string) {
	og_title       = extract_meta_property_content(page, "title")
	og_description = extract_meta_property_content(page, "description")
	og_site_name   = extract_meta_property_content(page, "site_name")
	page_title = ""

	return
}

func main() {
	url := "https://mateishome.page"
	body, _ := fetch_page(url)
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Printf("%s", body)
	// }
	page_title, og_title, og_description, og_site_name := extract_page_info(body)
	fmt.Println(page_title)
	fmt.Println(og_title)
	fmt.Println(og_description)
	fmt.Println(og_site_name)
}
