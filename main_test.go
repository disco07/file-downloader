package main

import (
	"errors"
	"fmt"
	"testing"
)

func TestDownloader(t *testing.T) {
	var downloaderTests = []struct {
		description string
		url         string
		expected    error
	}{
		{
			description: "Download work",
			url:         "https://agritrop.cirad.fr/584726/1/Rapport.pdf",
			expected:    nil,
		},
		{
			description: "invalid url",
			url:         "",
			expected:    fmt.Errorf("invalid url"),
		},
		{
			description: "unable to download file with multithreading",
			url:         "https://youtu.be/w0NQlEMjntI",
			expected:    errors.New("unable to download file with multithreads"),
		},
		{
			description: "unable to parse variable",
			url:         "https://github.com/disco07/file-downloader",
			expected:    errors.New("unable to parse variable"),
		},
	}

	for _, tt := range downloaderTests {
		t.Run(tt.description, func(t *testing.T) {
			err := worker(tt.url)
			if tt.expected == nil && err != nil {
				t.Errorf("Unexpected error for input %v: %v (expected %v)", tt.url, err, tt.expected)
			}

			if tt.expected != nil && err.Error() != tt.expected.Error() {
				t.Errorf("Unexpected error for input %v: %v (expected %v)", tt.url, err, tt.expected)
			}
		})
	}
}

func BenchmarkDownloader(b *testing.B) {
	for i := 0; i < b.N; i++ {
		worker("https://agritrop.cirad.fr/584726/1/Rapport.pdf")
	}
}
