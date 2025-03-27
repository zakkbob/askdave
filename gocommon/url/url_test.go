package url_test

import (
	"testing"

	"github.com/ZakkBob/AskDave/gocommon/url"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	parseTests := []struct {
		name string
		in   string
		out  string
		err  error
	}{
		{
			name: "https scheme",
			in:   "https://www.example.com",
			out:  "https://www.example.com",
			err:  nil,
		},
		{
			name: "http scheme",
			in:   "http://www.example.com",
			out:  "http://www.example.com",
			err:  nil,
		},
		{
			name: "fragment",
			in:   "https://www.example.com#fragment",
			out:  "https://www.example.com",
			err:  nil,
		},
		{
			name: "credentials",
			in:   "https://user_name:password@www.example.com",
			out:  "https://www.example.com",
			err:  nil,
		},
		{
			name: "invalid scheme",
			in:   "ftp://www.example.com",
			out:  "",
			err:  url.ErrInvalidScheme,
		},
		{
			name: "no scheme",
			in:   "www.example.com",
			out:  "www.example.com",
			err:  nil,
		},
		{
			name: "path",
			in:   "/example",
			out:  "/example",
			err:  nil,
		},
	}

	for _, pt := range parseTests {
		t.Run(pt.name, func(t *testing.T) {
			u, err := url.Parse(pt.in)

			if pt.err != nil {
				require.ErrorIs(t, err, pt.err, "Parse should return an error")
				return
			} else {
				require.NoError(t, err, "Parse should not return an error")
			}

			require.Equal(t, pt.out, u.String())
		})
	}
}

func TestParseAbs(t *testing.T) {
	parseTests := []struct {
		name string
		in   string
		out  string
		err  error
	}{
		{
			name: "https scheme",
			in:   "https://www.example.com",
			out:  "https://www.example.com",
			err:  nil,
		},
		{
			name: "http scheme",
			in:   "http://www.example.com",
			out:  "http://www.example.com",
			err:  nil,
		},
		{
			name: "fragment",
			in:   "https://www.example.com#fragment",
			out:  "https://www.example.com",
			err:  nil,
		},
		{
			name: "credentials",
			in:   "https://user_name:password@www.example.com",
			out:  "https://www.example.com",
			err:  nil,
		},
		{
			name: "invalid scheme",
			in:   "ftp://www.example.com",
			out:  "",
			err:  url.ErrInvalidScheme,
		},
		{
			name: "not absolute",
			in:   "www.example.com",
			out:  "",
			err:  url.ErrNotAbsolute,
		},
	}

	for _, pt := range parseTests {
		t.Run(pt.name, func(t *testing.T) {
			u, err := url.ParseAbs(pt.in)

			if pt.err != nil {
				require.ErrorIs(t, err, pt.err, "Parse should return an error")
				return
			} else {
				require.NoError(t, err, "Parse should not return an error")
			}

			require.Equal(t, pt.out, u.String())
		})
	}
}
