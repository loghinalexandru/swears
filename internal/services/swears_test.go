package services

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"os"
	"sync"
	"testing"

	"github.com/loghinalexandru/swears/internal/models"
)

func TestNew(t *testing.T) {
	t.Parallel()

	got := NewSwears([]models.SwearsRepo{TestRepo{}}, http.DefaultClient, "")

	if got.data != nil && got.data["en"] == nil {
		t.Error("different repositories")
	}

	if got.client != http.DefaultClient {
		t.Error("different http client")
	}
}

func TestGetSwear(t *testing.T) {
	t.Parallel()

	target := NewSwears([]models.SwearsRepo{TestRepo{}}, http.DefaultClient, "")
	got, err := target.GetSwear("en")

	if err != nil {
		t.Error(err)
	}

	if got == "" {
		t.Error("nil swear")
	}
}

func TestGetSwearWithInvalidLanguage(t *testing.T) {
	t.Parallel()

	target := NewSwears([]models.SwearsRepo{TestRepoWithError{}}, http.DefaultClient, "")
	_, err := target.GetSwear("en")

	if err == nil {
		t.Error("error is nill")
	}
}

func TestGetSwearWithRepoError(t *testing.T) {
	t.Parallel()

	target := NewSwears([]models.SwearsRepo{TestRepo{}}, http.DefaultClient, "")
	_, err := target.GetSwear("invalid language")

	if err == nil {
		t.Error("error is nill")
	}
}

func TestGetSwearFileWithInvalidLanguage(t *testing.T) {
	t.Parallel()

	target := NewSwears([]models.SwearsRepo{TestRepo{}}, http.DefaultClient, "")
	got := target.GetSwearFile("invalid language", false)

	if got != nil {
		t.Error("wrong value returned")
	}
}

func TestGetSwearFileWithError(t *testing.T) {
	t.Parallel()

	target := NewSwears([]models.SwearsRepo{TestRepoWithError{}}, http.DefaultClient, "")
	got := target.GetSwearFile("en", false)

	if got != nil {
		t.Error("wrong value returned")
	}
}

func TestGetSwearFilePlain(t *testing.T) {
	t.Parallel()

	client := newTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(`OK`)),
			Header:     make(http.Header),
		}
	})

	target := NewSwears([]models.SwearsRepo{TestRepo{}}, client, t.TempDir())
	got := target.GetSwearFile("en", true)

	if got == nil {
		t.Error("buffer is empty")
	}
}

func TestDownloadTTSFile(t *testing.T) {
	t.Parallel()

	client := newTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(`OK`)),
			Header:     make(http.Header),
		}
	})

	mock := Swears{
		client: client,
		mtx:    sync.Mutex{},
	}

	testFile := t.TempDir() + "/" + t.Name()
	err := mock.downloadTTSFile(testFile, "test_text", "test_lang")

	if err != nil {
		t.Fatal("This should be nil!")
	}

	data, err := os.ReadFile(testFile)

	if err != nil || string(data) != "OK" {
		t.Fatal("File not created")
	}
}

type TestRepo struct{}

func (t TestRepo) Get() (models.Record, error) {
	return models.Record{Value: "swear"}, nil
}

func (t TestRepo) Lang() string {
	return "en"
}

type TestRepoWithError struct{}

func (t TestRepoWithError) Get() (models.Record, error) {
	return models.Record{}, errors.New("test error")
}

func (t TestRepoWithError) Lang() string {
	return "en"
}

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func newTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}
