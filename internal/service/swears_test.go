package service

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"os"
	"sync"
	"testing"

	"github.com/loghinalexandru/swears/internal/model"
)

func TestNew(t *testing.T) {
	t.Parallel()

	got := NewSwears([]model.SwearsRepo{TestRepo{}}, "")

	if got.data != nil && got.data["en"] == nil {
		t.Error("missing expected \"en\" repository")
	}

	if got.client != http.DefaultClient {
		t.Error("unexpected http client set")
	}
}

func TestGetSwear(t *testing.T) {
	t.Parallel()

	target := NewSwears([]model.SwearsRepo{TestRepo{}}, "")
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

	target := NewSwears([]model.SwearsRepo{TestRepo{}}, "")
	_, err := target.GetSwear("invalid language")

	if err != errMissingRepo {
		t.Errorf("unexpected error, got: %q want: %q", err, errMissingRepo)
	}
}

func TestGetSwearFileWithInvalidLanguage(t *testing.T) {
	t.Parallel()

	target := NewSwears([]model.SwearsRepo{TestRepo{}}, "")
	_, err := target.GetSwearFile("invalid language", nil)

	if err != errMissingRepo {
		t.Errorf("unexpected error, got: %q want: %q", err, errMissingRepo)
	}
}

func TestGetSwearWithRepoError(t *testing.T) {
	t.Parallel()

	target := NewSwears([]model.SwearsRepo{TestRepoWithError{}}, "")
	_, err := target.GetSwear("en")

	if err == nil {
		t.Error("missing expected error")
	}
}

func TestGetSwearFileWithRepoError(t *testing.T) {
	t.Parallel()

	target := NewSwears([]model.SwearsRepo{TestRepoWithError{}}, "")
	_, err := target.GetSwearFile("en", nil)

	if err == nil {
		t.Error("missing expected error")
	}
}

func TestGetSwearFilePlain(t *testing.T) {
	t.Parallel()

	testRepos := []model.SwearsRepo{TestRepo{}}
	tempDirPath := t.TempDir()
	swearRecord, _ := testRepos[0].Get()

	client := newTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(`OK`)),
			Header:     make(http.Header),
		}
	})

	target := NewSwears(testRepos, tempDirPath, WithClient(client))
	got, err := target.GetSwearFile("en", nil)

	if err != nil {
		t.Fatalf("unexpected error: %q", err)
	}

	if got == nil {
		t.Fatal("resulting buffer is empty")
	}

	if _, err := os.Stat(tempDirPath + "/" + swearRecord.ID.String() + ".mp3"); err != nil {
		t.Error(err)
	}
}

func TestGetSwearFileEncoded(t *testing.T) {
	t.Parallel()

	want := "bogus encoded string"
	testRepos := []model.SwearsRepo{TestRepo{}}
	tempDirPath := t.TempDir()

	client := newTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(`OK`)),
			Header:     make(http.Header),
		}
	})

	target := NewSwears(testRepos, tempDirPath, WithClient(client))
	got, err := target.GetSwearFile("en", TestEncoder{
		msg: want,
	})

	if err != nil {
		t.Fatalf("unexpected error: %q", err)
	}

	if got == nil {
		t.Fatal("encoded buffer is empty")
	}

	if string(got) != want {
		t.Error("invalid encoded buffer")
	}
}

func TestGetSwearFileEncodedWithErr(t *testing.T) {
	t.Parallel()

	testRepos := []model.SwearsRepo{TestRepo{}}
	tempDirPath := t.TempDir()

	client := newTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(`OK`)),
			Header:     make(http.Header),
		}
	})

	target := NewSwears(testRepos, tempDirPath, WithClient(client))
	_, err := target.GetSwearFile("en", TestEncoderWithErr{})

	if err == nil {
		t.Error("missing expected error")
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
		t.Fatalf("unexpected error : %q", err)
	}

	data, err := os.ReadFile(testFile)

	if err != nil || string(data) != "OK" {
		t.Fatal("file not created")
	}
}

type TestRepo struct{}

type TestEncoder struct {
	msg string
}

type TestEncoderWithErr struct {
}

func (t TestEncoder) Encode(io.Reader) ([]byte, error) {
	return []byte(t.msg), nil
}

func (t TestEncoderWithErr) Encode(io.Reader) ([]byte, error) {
	return nil, errors.New("bogus error")
}
func (t TestRepo) Get() (model.Record, error) {
	return model.Record{Value: "swear"}, nil
}

func (t TestRepo) Lang() string {
	return "en"
}

type TestRepoWithError struct{}

func (t TestRepoWithError) Get() (model.Record, error) {
	return model.Record{}, errors.New("test error")
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
