package rsp

import (
	"fmt"
	"testing"
)

func TestCreateShortURL(t *testing.T) {
	type testpair struct {
		longURL  string // long url to be shortened
		code     int    // status code returned from remote api
		shortURL string // short url returned from remote api
	}

	// This is the test to verify that we get the 'right' result from external api
	// with the assumption that bitly caches the short urls for a long time.
	// And google url is chosen because it is a highly active website and
	// less likely to expire.
	var tests = []testpair{
		{"https://www.google.com", 200, "http://bit.ly/2SmO2Qr"},
		{"http://invalidurl", 400, ""},
	}

	for _, pair := range tests {
		t.Run(fmt.Sprintf("shorten %s", pair.longURL), func(t *testing.T) {
			u, c := CreateShortURL(pair.longURL)
			if u != pair.shortURL || c != pair.code {
				t.Error(
					"For", pair.longURL,
					"expected", pair.shortURL, pair.code,
					"got", u, c,
				)
			}
		})
	}
}
