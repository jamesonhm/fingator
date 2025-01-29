package openfigi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jamesonhm/fingator/internal/openfigi/models"
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

func (c *Client) CallURL(ctx context.Context, uri string, params, response any) error {
	c.limiter.Wait(ctx)
	uri = c.baseurl + uri

	body, err := json.Marshal(params)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("X-OPENFIGI-APIKEY", c.apiKey)
	resp, err := c.httpC.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	fmt.Println("resp Body:", resp.Body)
	fmt.Println("resp Code:", resp.Status)
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return fmt.Errorf("error decoding json: %w", err)
	}

	return nil
}

func (c *Client) Mapping(ctx context.Context, params []models.MappingRequest) (*[]models.MappingResponse, error) {
	res := &[]models.MappingResponse{}
	err := c.CallURL(ctx, mappingURL, params, res)
	return res, err
}
