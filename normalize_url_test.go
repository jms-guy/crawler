package main

import (
	"net/url"
	"testing"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name string
		inputURL string
		expected string
	}{
		{
			name: "remove scheme",
			inputURL: "https://blog.boot.dev/path",
			expected: "http://blog.boot.dev/path",
		},
		{
			name: "remove trailing /",
			inputURL: "blog.boot.dev/path/",
			expected: "http://blog.boot.dev/path",
		},
		{
			name: "remove scheme and trailing",
			inputURL: "http://blog.boot.dev/path/",
			expected: "http://blog.boot.dev/path",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			baseURL, _ := url.Parse(tc.inputURL)
			actual, err := normalizeURL(baseURL)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			}
			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}