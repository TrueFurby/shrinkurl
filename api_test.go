package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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

	resp, err := http.PostForm(apiurl, payload)
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

func TestCheck(t *testing.T) {
	var (
		hash   = makeHash("google.com")
		query  = url.Values{"hash": []string{hash}}
		apiurl = server.URL + "/check" + fmt.Sprintf("?%s", query.Encode())
	)

	dburl := Url{
		Id:   1,
		Hash: hash,
		Url:  "http://google.com",
	}
	mockStore.urls[hash] = dburl

	resp, err := http.Get(apiurl)
	if err != nil {
		t.Error(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("response: %s", string(body))
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected 200 but got: %d", resp.StatusCode)
	}

	b, _ := json.Marshal(dburl)
	if !bytes.Equal(b, bytes.TrimSpace(body)) {
		t.Errorf("expected '%s' but got: '%s'", string(b), string(body))
	}
}
