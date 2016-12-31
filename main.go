package main

import (
	"errors"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	url                *url.URL
	httpClient         *http.Client
	username, password string
}

func newClient(username, password string) (*Client, error) {
	u, _ := url.Parse("https://gateway.watsonplatform.net/language-translation/api/v2")

	hc := &http.Client{
		Timeout: 5 * time.Second,
	}

	if len(username) == 0 {
		return nil, errors.New("missing username")
	}

	if len(password) == 0 {
		return nil, errors.New("missing user password")
	}

	return &Client{
		url:        u,
		httpClient: hc,
		username:   username,
		password:   password,
	}, nil
}
