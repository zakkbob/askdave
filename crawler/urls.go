package main

//import "fmt"

type url struct {
	protocol      string
	subdomain     string
	domain        string
	tld           string
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
