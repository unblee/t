package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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
	option:
		-v --version: show version

	t translates input text specified by argument or STDIN using Watson Language Translation API.
	Source language will be automatically detected.

	export T_WATSON_LANGUAGE_TRANSLATOR_API_USERNAME=<Your Watson Language Translator API username>
	export T_WATSON_LANGUAGE_TRANSLATOR_API_PASSWORD=<Your Watson Language Translator API password>

	Example:
		$ t Good morning!
		おはようございます!
		$ t おはようございます!
		Good morning!
`

var (
	Version   string
	Revision  string
	GoVersion string
	ShowVerL  = flag.Bool("version", false, "show version")
	ShowVerS  = flag.Bool("v", false, "show version")
)

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
	if *ShowVerL || *ShowVerS {
		fmt.Printf("version:    %s\n", Version)
		fmt.Printf("revision:   %s\n", Revision)
		fmt.Printf("go version: %s\n", GoVersion)
		os.Exit(ExitCodeOK)
	}
	os.Exit(run(flag.Args()))
}

func run(args []string) int {
	logger := log.New(os.Stderr, "[ERROR] ", log.LstdFlags)

	var text string
	if isatty.IsTerminal(os.Stdin.Fd()) { // with Args
		if flag.NArg() == 0 {
			fmt.Println(usageMsg)
			os.Exit(ExitCodeOK)
		}
		text = strings.Join(args, " ")
	} else { // with Pipe
		b, _ := ioutil.ReadAll(os.Stdin)
		text = string(b)
	}

	reqBody := newRequestBody(text)

	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		logger.Println(err)
		return ExitCodeError
	}

	username := os.Getenv("T_WATSON_LANGUAGE_TRANSLATOR_API_USERNAME")
	password := os.Getenv("T_WATSON_LANGUAGE_TRANSLATOR_API_PASSWORD")
	client, err := newClient(username, password)
	if err != nil {
		logger.Println(err)
		return ExitCodeError
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := client.translate(ctx, bytes.NewReader(reqJSON))
	if err != nil {
		logger.Println(err)
		return ExitCodeError
	}

	fmt.Print(result)

	return ExitCodeOK
}

type Client struct {
	url                *url.URL
	httpClient         *http.Client
	username, password string
}

func newClient(username, password string) (*Client, error) {
	u, _ := url.Parse("https://gateway.watsonplatform.net/language-translator/api/v2")

	hc := new(http.Client)

	if len(username) == 0 {
		return nil, errors.New("Missing username. Please set 'T_WATSON_LANGUAGE_TRANSLATOR_API_USERNAME'.")
	}

	if len(password) == 0 {
		return nil, errors.New("Missing user password. Please set 'T_WATSON_LANGUAGE_TRANSLATOR_API_PASSWORD'.")
	}

	return &Client{
		url:        u,
		httpClient: hc,
		username:   username,
		password:   password,
	}, nil
}

func (c *Client) translate(ctx context.Context, body io.Reader) (string, error) {
	req, err := c.newRequest(ctx, "POST", "/translate", body)
	if err != nil {
		return "", err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		errBody := new(errorResponseBody)
		decodeBody(resp, errBody)

		code := resp.StatusCode
		switch code {
		case 400, 404, 500:
			return "", fmt.Errorf("\"code\":%d, \"error_message\":\"%s\"", code, errBody.ErrorMessage)
		default:
			return "", fmt.Errorf("\"code\":%d, \"error\":\"%s\", \"description\":\"%s\"", resp.StatusCode, errBody.Error, errBody.Description)
		}
	}

	respBody := new(responseBody)
	decodeBody(resp, respBody)
	return respBody.Translations[0].Translation, nil
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

func decodeBody(resp *http.Response, out interface{}) error {
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(out)
}

type requestBody struct {
	ModelID string   `json:"model_id"`
	Source  string   `json:"source"`
	Target  string   `json:"target"`
	Text    []string `json:"text"`
}

func newRequestBody(text string) *requestBody {
	model, source, target := detectLang(text)
	return &requestBody{
		ModelID: model,
		Source:  source,
		Target:  target,
		Text:    []string{text},
	}
}

func detectLang(text string) (model, source, target string) {
	for n, c := range text {
		if n > 20 {
			break
		}
		if c > 127 {
			return "ja-en", "ja", "en"
		}
	}
	return "en-ja", "en", "ja"
}

type responseBody struct {
	Translations []struct {
		Translation string `json:"translation"`
	} `json:"translations"`
}

type errorResponseBody struct {
	ErrorMessage string `json:"error_message"`
	Error        string `json:"error"`
	Description  string `json:"description"`
}
