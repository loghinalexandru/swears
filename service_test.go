package main

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"sync"
	"testing"
)

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestDownloadTTSFile(t *testing.T) {
	client := NewTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(`OK`)),
			Header:     make(http.Header),
		}
	})

	mock := SwearsSvc{
		client: client,
		mtx:    sync.Mutex{},
	}

	testFile := "test.mp3"
	err := mock.downloadTTSFile(testFile, "test_text", "test_lang")
	defer os.Remove(testFile)

	if err != nil {
		t.Error("This should be nil!")
		t.FailNow()
	}

	data, err := os.ReadFile(testFile)

	if err != nil || string(data) != "OK" {
		t.Error("File not created")
	}
}
