package main

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"path"
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

func (c *Client) newRequest(ctx context.Context, method, spath string, body io.Reader) (*http.Request, error) {
	u := *c.url
	u.Path = path.Join(c.url.Path, spath)

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return req, nil
}

type requestBody struct {
	ModelID string   `json:"model_id"`
	Source  string   `json:"source"`
	Target  string   `json:"target"`
	Text    []string `json:"text"`
}
