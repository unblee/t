package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/mattn/go-isatty"
)

const usageMsg = `
Usage: t [input text] or echo [input text] | t

	t translates input text specified by argument or STDIN using Watson Language Translation API.
	Source language will be automatically detected.

	export T_WATSON_LANGUAGE_TRANSLATOR_API_USERNAME = <Your Watson Language Translator API username>
	export T_WATSON_LANGUAGE_TRANSLATOR_API_PASSWORD = <Your Watson Language Translator API password>

	Example:
		$ t Good morning!
		おはようございます!
		$ t おはようございます!
		Good morning!
`

const (
	ExitCodeOK = iota
	ExitCodeError
)

func main() {
	flag.Usage = func() {
		fmt.Println(usageMsg)
		os.Exit(ExitCodeOK)
	}
	flag.Parse()
	os.Exit(run(flag.Args()))
}

func run(args []string) int {
	var text string
	if isatty.IsTerminal(os.Stdin.Fd()) {
		if flag.NArg() == 0 {
			fmt.Println(usageMsg)
			os.Exit(ExitCodeOK)
		}
		text = strings.Join(args, " ")
	} else { // with Pipe
		b, _ := ioutil.ReadAll(os.Stdin)
		text = string(b)
	}
	return ExitCodeOK
}

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
