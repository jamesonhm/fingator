package polygon

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jamesonhm/fingator/internal/uri"
)

const (
	APIURL = "https://api.polygon.io"
)

type Client struct {
	baseurl    string
	apiKey     string
	httpC      http.Client
	uriBuilder *uri.URIBuilder
}

func New(apiKey string, timeout time.Duration) Client {
	return Client{
		baseurl: APIURL,
		apiKey:  apiKey,
		httpC: http.Client{
			Timeout: timeout,
		},
		uriBuilder: uri.New(APIURL),
	}
}

// Call makes API call based on path and params
func (c *Client) Call(ctx context.Context, path string, params, response any) error {
	uri := c.uriBuilder.EncodeParams(path, params)
	//if err != nil {
	//	return err
	//}
	fmt.Printf("clietn-call-uri: %s\n", uri)
	return c.CallURL(ctx, uri, response)
}

func (c *Client) CallURL(ctx context.Context, uri string, response any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+c.apiKey)
	resp, err := c.httpC.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return fmt.Errorf("error decoding json: %w", err)
	}

	return nil
}
