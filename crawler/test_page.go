package main

import (
	"testing"
)

func TestAddLink(t *testing.T) {
	var p page
	expected := "https://google.com"
	u, _ := parseAbsoluteUrl(expected)
	p.addLink(u)
	got := p.links[0].String()
	if got != expected{
		t.Errorf("got '%s', expected '%s'", got, expected)
	}
	
}