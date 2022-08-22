package main

import (
	"fmt"
	"testing"
)

var downloaderTests = []struct {
	url      string
	expected error
}{
	{"https://agritrop.cirad.fr/584726/1/Rapport.pdf", nil},
	{"", fmt.Errorf("invalid url")},
}

func TestDownloader(t *testing.T) {
	for _, tt := range downloaderTests {
		err := downloader(tt.url)
		if err == nil {
			t.Errorf("Expected error for input %v", tt.url)
		}
		if err.Error() != tt.expected.Error() {
			t.Errorf("Unexpected error for input %v: %v (expected %v)", tt.url, err, tt.expected)
		}
	}
}

func BenchmarkDownloader(b *testing.B) {
	for i := 0; i < b.N; i++ {
		downloader("https://agritrop.cirad.fr/584726/1/Rapport.pdf")
	}
}
