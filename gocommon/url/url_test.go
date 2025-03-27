package url_test

import (
	"testing"

	"errors"

	"github.com/ZakkBob/AskDave/gocommon/url"
	"github.com/stretchr/testify/require"
)

type parseTest struct {
	in          string
	out         string
	parseErr    error
	parseAbsErr error
}

var parseTests = []parseTest{
	{
		in:          "https://www.example.com",
		out:         "https://www.example.com",
		parseErr:    nil,
		parseAbsErr: nil,
	},
	{
		in:          "http://www.example.com",
		out:         "http://www.example.com",
		parseErr:    nil,
		parseAbsErr: nil,
	},
	{
		in:          "https://www.example.com#fragment",
		out:         "https://www.example.com",
		parseAbsErr: nil,
		parseErr:    nil,
	},
	{
		in:          "https://user_name:password@www.example.com",
		out:         "https://www.example.com",
		parseErr:    nil,
		parseAbsErr: nil,
	},
	{
		in:          "ftp://www.example.com",
		out:         "",
		parseErr:    url.ErrInvalidScheme,
		parseAbsErr: url.ErrInvalidScheme,
	},
	{
		in:          "www.example.com",
		out:         "www.example.com",
		parseErr:    nil,
		parseAbsErr: nil,
	},
	{
		in:          "/example",
		out:         "/example",
		parseErr:    nil,
		parseAbsErr: url.ErrNotAbsolute,
	},
}

func TestParse(t *testing.T) {
	for _, test := range parseTests {
		u, err := url.Parse(test.in)

		if !errors.Is(err, test.parseAbsErr) {
			t.Fatalf("expected '%v' error with url %s", test.parseAbsErr, test.in)
		}

		if err != nil {
			return
		}

		require.Equal(t, test.out, u.String())
	}
}

func TestParseAbs(t *testing.T) {
	for _, test := range parseTests {
		u, err := url.ParseAbs(test.in)

		if !errors.Is(err, test.parseAbsErr) {
			t.Fatalf("expected '%v' error with url %s", test.parseAbsErr, test.in)
		}

		if err != nil {
			return
		}

		require.Equal(t, test.out, u.String())
	}
}

type stringNoPathTest struct {
	in  string
	out string
}

var stringNoPathTests = []stringNoPathTest{
	{
		in:  "https://example.com/path/example/test/",
		out: "https://example.com",
	},
}

func TestStringNoPath(t *testing.T) {
	for _, test := range stringNoPathTests {
		u, err := url.Parse(test.in)
		require.NoError(t, err, "Parse should not return an error")

		require.Equal(t, test.out, u.StringNoPath())
	}
}
