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
		t.Error("different repositories")
	}

	if got.client != http.DefaultClient {
		t.Error("different http client")
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

	target := NewSwears([]model.SwearsRepo{TestRepoWithError{}}, "")
	_, err := target.GetSwear("en")

	if err == nil {
		t.Error("error is nill")
	}
}

func TestGetSwearWithRepoError(t *testing.T) {
	t.Parallel()

	target := NewSwears([]model.SwearsRepo{TestRepo{}}, "")
	_, err := target.GetSwear("invalid language")

	if err == nil {
		t.Error("error is nill")
	}
}

func TestGetSwearFileWithInvalidLanguage(t *testing.T) {
	t.Parallel()

	target := NewSwears([]model.SwearsRepo{TestRepo{}}, "")
	got := target.GetSwearFile("invalid language", nil)

	if got != nil {
		t.Error("wrong value returned")
	}
}

func TestGetSwearFileWithError(t *testing.T) {
	t.Parallel()

	target := NewSwears([]model.SwearsRepo{TestRepoWithError{}}, "")
	got := target.GetSwearFile("en", nil)

	if got != nil {
		t.Error("wrong value returned")
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
	got := target.GetSwearFile("en", nil)

	if got == nil {
		t.Fatal("buffer is empty")
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
	got := target.GetSwearFile("en", TestEncoder{
		msg: want,
	})

	if got == nil {
		t.Fatal("buffer is empty")
	}

	if string(got) != want {
		t.Error("invalid encoded stream")
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
	got := target.GetSwearFile("en", TestEncoderWithErr{})

	if got != nil {
		t.Error("buffer is not empty")
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
		t.Fatal("this should be nil")
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
