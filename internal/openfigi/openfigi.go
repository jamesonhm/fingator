package openfigi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

const (
	baseURL    = "https://api.openfigi.com"
	mappingURL = "/v3/mapping"
)

type Client struct {
	baseurl string
	apiKey  string
	httpC   http.Client
	limiter *rate.Limiter
}

func New(apiKey string, timeout time.Duration, reqsPerSec rate.Limit) Client {
	return Client{
		baseurl: baseURL,
		apiKey:  apiKey,
		httpC: http.Client{
			Timeout: timeout,
		},
		limiter: rate.NewLimiter(reqsPerSec, 1),
	}
}

func (c *Client) CallURL(ctx context.Context, uri string, response any) error {
	c.limiter.Wait(ctx)
	uri = c.baseurl + uri
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return err
	}

	req.Header.Add("X-OPENFIGI-APIKEY", c.apiKey)
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

func (c *Client) MapCUSIPs(ctx context.Context, cusips []string) (*models.GroupedDailyResponse, error) {
	res := &models.GroupedDailyResponse{}
	err := c.CallURL(ctx, GroupedDailyPath, params, res)
	return res, err
}
