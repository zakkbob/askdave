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

// Parse parses a [URL] in the context of the receiver.
// The provided URL may be relative or absolute.
// Parse returns nil, err on parse failure.
// Normalises a URL by stripping the query, fragment, and user.
// Returns an ErrInvalidScheme if the scheme is not 'http' or 'https'
// Returns ErrNotAbsolute if the scheme is empty
func (u *URL) Parse(ref string) (*URL, error) {
	parsed, err := u.URL.Parse(ref)
	if err != nil {
		return nil, err
	}

	newURL := &URL{URL: *parsed}
	err = normaliseURL(newURL)
	if err != nil {
		return nil, err
	}

	return newURL, err
}

// Returns url string without path
func (u *URL) StringNoPath() string {
	copy := *u
	copy.Path = ""
	return copy.String()
}

// Parses an absolute raw url into a [URL] struct
//
// Normalises a URL by stripping the query, fragment, and user.
// Returns an ErrInvalidScheme if the scheme is not 'http' or 'https'
// Returns ErrNotAbsolute if the scheme is empty
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

// Normalises a URL by stripping the query, fragment, and user.
// Returns an ErrInvalidScheme if the scheme is not ”, 'http' or 'https'
func normaliseURL(u *URL) error {
	if u.Scheme != "http" && u.Scheme != "https" && u.Scheme != "" {
		return ErrInvalidScheme
	}

	u.RawQuery = ""
	u.Fragment = ""
	u.RawFragment = ""
	u.User = nil

	return nil
}

// Parses a raw url into a [URL] struct
//
// Normalises a URL by stripping the query, fragment, and user.
// Returns an ErrInvalidScheme if the scheme is not ”, 'http' or 'https'
func Parse(rawURL string) (*URL, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse raw url: %w", err)
	}

	newURL := &URL{URL: *u}

	err = normaliseURL(newURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse raw url: %w", err)
	}

	return newURL, nil
}

// Parses an array of strings into a [URL] array
//
// Normalises a URL by stripping the query, fragment, and user.
// Returns an ErrInvalidScheme if the scheme is not ”, 'http' or 'https'
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
	tmp, err := Parse(s)
	if err != nil {
		return fmt.Errorf("failed to parse url: %w", err)
	}
	*u = *tmp
	return nil
}
