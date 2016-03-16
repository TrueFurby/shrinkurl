package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var (
	mockStore = NewMockStorage()
	server    *httptest.Server
)

func init() {
	var app = &App{
		Storage: mockStore,
		Cache:   NewCache(10),
	}
	app.Router = NewRouter(app)
	server = httptest.NewServer(app.Router)
}

func TestAdd(t *testing.T) {
	var (
		apiurl  = server.URL + "/add"
		payload = url.Values{"url": []string{"google.com"}}
	)

	resp, err := http.DefaultClient.PostForm(apiurl, payload)
	if err != nil {
		t.Error(err)
	}

	defer resp.Body.Close()
	if body, err := ioutil.ReadAll(resp.Body); err != nil {
		t.Error(err)
	} else {
		t.Logf("response: %s", string(body))
	}

	if resp.StatusCode != 201 {
		t.Errorf("expected 201 but got: %d", resp.StatusCode)
	}
}
