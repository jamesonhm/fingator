package polygon

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-playground/form/v4"
)

const (
	APIURL = "https://api.polygon.io"
)

type Client struct {
	baseurl      string
	apiKey       string
	httpC        http.Client
	pathEncoder  *form.Encoder
	queryEncoder *form.Encoder
}

func New(apiKey string, timeout time.Duration) Client {
	return Client{
		baseurl: APIURL,
		apiKey:  apiKey,
		httpC: http.Client{
			Timeout: timeout,
		},
		pathEncoder:  newEncoder("path"),
		queryEncoder: newEncoder("query"),
	}
}

func (c *Client) Call(ctx context.Context, path string, params, response any) error {
	// TODO: add baseurl here?
	uri, err := c.encodeParams(path, params)
	if err != nil {
		return err
	}
	fmt.Println(uri)
	return nil
}

func (c *Client) encodeParams(path string, params any) (string, error) {
	vals, err := c.pathEncoder.Encode(&params)
	if err != nil {
		return "", fmt.Errorf("error encoding path params: %w", err)
	}
	fmt.Printf("encoded vals: %v\n", vals)

	pathParams := map[string]string{}
	for k, v := range vals {
		pathParams[k] = v[0]
	}
	fmt.Printf("map pathParams: %v\n", pathParams)
	for k, v := range pathParams {
		path = strings.ReplaceAll(path, fmt.Sprintf("{%s}", k), url.PathEscape(v))
	}
	return path, nil
}

func newEncoder(tag string) *form.Encoder {
	e := form.NewEncoder()
	e.SetMode(form.ModeExplicit)
	e.SetTagName(tag)

	return e
}
