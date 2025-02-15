package main

import (
	"fmt"
	"strconv"
)

func main() {
	var r taskRunner

	for i := 0; i < 25; i++ {
		s := "https://example.com/robots.txt" + strconv.Itoa(i)

		u, _ := parseAbsoluteUrl(s)
		r.t.robots.slice = append(r.t.robots.slice, u)
	}

	for i := 0; i < 25; i++ {
		u, _ := parseAbsoluteUrl(fmt.Sprintf("https://example.com/sitemap%d.xml", i))
		r.t.sitemaps.slice = append(r.t.sitemaps.slice, u)
	}

	for i := 0; i < 25; i++ {
		u, _ := parseAbsoluteUrl(fmt.Sprintf("https://example.com/page%d.html", i))
		r.t.pages.slice = append(r.t.pages.slice, u)
	}

	r.run(25)
}
