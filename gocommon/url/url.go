//----------------------------------------------------------------//
// Thin abstraction over net/url which provides some helper funcs //
//----------------------------------------------------------------//

package url

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
)

type URL struct {
	url.URL
}

var ErrInvalidScheme = errors.New("url does not contain a valid scheme")
var ErrNotAbsolute = errors.New("url is not absolute")

func (u *URL) Parse(ref string) (*URL, error) {
	parsed, err := u.URL.Parse(ref)
	return &URL{URL: *parsed}, err
}

// Returns url string without path
func (u *URL) StringNoPath() string {
	copy := *u
	copy.Path = ""
	return copy.String()
}

// Parses an absolute raw url into a [URL] struct
//
// Normalises the url by removing the fragment and credentials.
// Will return an error if the scheme is not 'http' or 'https'.
// Will return an error if the url is not absolute.
func ParseAbs(rawURL string) (*URL, error) {
	u, err := Parse(rawURL)
	if err != nil {
		return nil, err
	}

	if !u.IsAbs() {
		return nil, ErrNotAbsolute
	}

	return u, nil
}

// Parses a raw url into a [URL] struct
//
// Normalises the url by removing the fragment and credentials.
// Will return an error if the scheme is not blank or 'http' or 'https'.
func Parse(rawURL string) (*URL, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse raw url: %w", err)
	}

	if u.Scheme != "http" && u.Scheme != "https" && u.Scheme != "" {
		return nil, ErrInvalidScheme
	}

	u.RawQuery = ""
	u.Fragment = ""
	u.RawFragment = ""
	u.User = nil

	return &URL{URL: *u}, nil
}

// Parses an array of strings into a [URL] array
//
// Uses ParseAbs
func ParseMany(rawURLs []string) ([]*URL, error) {
	urls := make([]*URL, 0)
	var u *URL
	var err error
	for _, urlS := range rawURLs {
		u, err = ParseAbs(urlS)
		if err != nil {
			return nil, err
		}
		urls = append(urls, u)
	}
	return urls, nil
}

func (u *URL) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}

func (u *URL) UnmarshalJSON(data []byte) error {
	var s string
	json.Unmarshal(data, &s)
	tmp, err := url.Parse(s)
	if err != nil {
		return fmt.Errorf("failed to parse url: %w", err)
	}
	*u = URL{URL: *tmp}
	return nil
}
