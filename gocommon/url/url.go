//----------------------------------------------------------------//
// Thin abstraction over net/url which provides some helper funcs //
//----------------------------------------------------------------//

package url

import (
	"errors"
	"fmt"
	"net/url"
)

type URL = url.URL

var ErrInvalidScheme = errors.New("url does not contain a valid scheme")
var ErrNotAbsolute = errors.New("url is not absolute")

// Parses an absolute raw url into a [URL] struct
//
// Normalises the url by removing the fragment and credentials.
// Will return an error if the scheme is not 'http' or 'https'.
// Will return an error if the url is not absolute.
func ParseAbs(rawURL string) (*URL, error) {
	u, err := Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse absolute url: %w", err)
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

	u.Fragment = ""
	u.RawFragment = ""
	u.User = nil

	return u, nil
}
