package openfigi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jamesonhm/fingator/internal/openfigi/models"
	"github.com/jamesonhm/fingator/internal/rate"
)

const (
	baseURL    = "https://api.openfigi.com"
	mappingURL = "/v3/mapping"
)

type Client struct {
	baseurl   string
	apiKey    string
	httpC     http.Client
	limiter   *rate.Limiter
	Batchsize int
}

func New(apiKey string, timeout time.Duration) Client {
	var period time.Duration
	var batchsize int
	if apiKey == "" {
		period = time.Minute * 1
		batchsize = 10
	} else {
		period = time.Second * 6
		batchsize = 100
	}
	count := 25
	return Client{
		baseurl: baseURL,
		apiKey:  apiKey,
		httpC: http.Client{
			Timeout: timeout,
		},
		limiter:   rate.New(period, count),
		Batchsize: batchsize,
	}
}

func (c *Client) CallURL(ctx context.Context, uri string, params, response any) error {
	err := c.limiter.Wait(ctx)
	if err != nil {
		return err
	}
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
