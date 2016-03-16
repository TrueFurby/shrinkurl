package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"net/url"
)

type UrlStore interface {
	Update(u *Url) error
	Get(hash string) (Url, error)
	GetByUrl(url string) (Url, error)
	Remove(hash string) error
}

type Url struct {
	Id   int64  `json:"id"`
	Hash string `json:"hash"`
	Url  string `json:"url"`
}

func makeHash(u string) string {
	data := []byte(u)
	hash := fmt.Sprintf("%x", md5.Sum(data))
	return hash[:6]
}

func parseUrl(input string) (*url.URL, error) {
	if input == "" {
		return nil, errors.New("empty url entered")
	}
	u, err := url.Parse(input)
	if err != nil {
		return nil, err
	}

	if u.Scheme == "" {
		if u, err = url.Parse("http://" + input); err != nil {
			return nil, err
		}
	} else if u.Scheme != "http" && u.Scheme != "https" {
		return nil, errors.New("only http(s) scheme supported")
	}

	return u, nil
}
