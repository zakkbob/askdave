package main

import (
	"fmt"
	"regexp"
	"strings"
)

type urlValidator struct {
	allowedUrlPatterns   []regexp.Regexp
	diallowedUrlPatterns []regexp.Regexp
}

func extractUserAgentBlocks(content string) (blocks map[string]string) {
	blocks = make(map[string]string)

	agentRegex := regexp.MustCompile("(?si)User-Agent:")
	matches := agentRegex.Split(content, -1)

	fmt.Println(matches)

	blockRegex := regexp.MustCompile("(?si)(.*?)(?:\r\n|\r|\n)(.*)") //Dodgy regex v2 (yes i know i can do this without regex, but regex = high iq)
	
	for _, match := range matches {
		captures := blockRegex.FindStringSubmatch(match)
		
		if len(captures) < 3 {
			continue
		}

		userAgent := strings.TrimSpace(captures[1])
		directives := strings.TrimSpace(captures[2])

		blocks[userAgent] = directives
	}

	return blocks
}

func processRobotsTxt(content string) (validator urlValidator, sitemapUrl string) {
	validator = urlValidator {make([]regexp.Regexp, 0), make([]regexp.Regexp, 0) }
	
	//agentRegex := regexp.MustCompile("(?i)User-Agent:(.*)")
	//disallowRegex := regexp.MustCompile("(?i)Disallow:(.*)")
	//allowRegex := regexp.MustCompile("(?i)Allow:(.*)")
	
	return validator, ""
}
