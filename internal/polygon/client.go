package polygon

import (
	"context"
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
	fmt.Println(uri)
	return nil
}
