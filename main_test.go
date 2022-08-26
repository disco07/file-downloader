package main

import (
	"errors"
	"fmt"
	"testing"
)

var downloaderTestsNotError = []struct {
	url      string
	expected error
}{
	{"https://agritrop.cirad.fr/584726/1/Rapport.pdf", nil},
}

var downloaderTests = []struct {
	url      string
	expected error
}{
	{"", fmt.Errorf("invalid url")},
	{"https://youtu.be/w0NQlEMjntI", errors.New("unable to download file with multithreads")},
	{"https://github.com/disco07/file-downloader", errors.New("unable to parse variable")},
}

func TestDownloaderNotError(t *testing.T) {
	for _, tt := range downloaderTestsNotError {
		err := worker(tt.url)
		if err != nil {
			t.Fatalf("Expected error for input %v", tt.url)
		}
	}
}

func TestDownloader(t *testing.T) {
	for _, tt := range downloaderTests {
		err := worker(tt.url)
		if err.Error() != tt.expected.Error() {
			t.Errorf("Unexpected error for input %v: %v (expected %v)", tt.url, err, tt.expected)
		}
	}
}

func BenchmarkDownloader(b *testing.B) {
	for i := 0; i < b.N; i++ {
		worker("https://agritrop.cirad.fr/584726/1/Rapport.pdf")
	}
}
