package main

import (
	"strconv"
)

type url struct {
	protocol      string
	subdomain     string
	domain        string
	tld           string
	port          int
	path          []string
	trailingSlash bool
}

func (u *url) String() string {
	s := ""

	if u.protocol != "" {
		s += u.protocol + "://"
	}

	if u.subdomain != "" {
		s += u.subdomain + "."
	}

	s += u.domain + "."
	s += u.tld

	if (u.port != 0) && !(u.protocol == "http" && u.port == 80) && !(u.protocol == "https" && u.port == 443) {
		s += ":" + strconv.Itoa(u.port)
	}

	if len(u.path) != 0 {
		for _, p := range u.path {
			s += "/" + p
		}
	}

	if u.trailingSlash {
		s += "/"
	}

	return s
}
