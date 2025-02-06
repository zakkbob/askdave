package main

import (
	"fmt"
	"testing"
)

func TestExtractUserAgentBlocks(t *testing.T) {
	robotsTxt := "User-agent: * \nDisallow: / \n\nUser-agent: Googlebot \nDisallow: \n\nUser-agent: bingbot \nDisallow: /not-for-bing/ "

	blocks := extractUserAgentBlocks(robotsTxt)

	//t.Log(blocks)
	fmt.Println(blocks)
}