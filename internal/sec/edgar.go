package edgar

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jamesonhm/fingator/internal/sec/models"
	"github.com/jamesonhm/fingator/internal/uri"
)

const (
	CompanyTickerCIKPath = "https://www.sec.gov/files/company_tickers_exchange.json"
	BaseURL              = "https://data.sec.gov"
	CompanyFactsPath     = "/api/xbrl/companyfacts/CIK{cik_padded}.json"
)

type Client struct {
	agentName  string
	agentEmail string
	baseurl    string
	httpC      http.Client
	uriBuilder *uri.URIBuilder
}

func New(agentName, agentEmail string, timeout time.Duration) Client {
	return Client{
		agentName:  agentName,
		agentEmail: agentEmail,
		baseurl:    BaseURL,
		httpC: http.Client{
			Timeout: timeout,
		},
		uriBuilder: uri.New(),
	}
}

// Call makes API call based on path and params
func (c *Client) Call(ctx context.Context, path string, params, response any) error {
	uri := c.uriBuilder.EncodeParams(path, params)
	uri = c.baseurl + uri
	fmt.Printf("client-call-uri: %s\n", uri)
	return c.CallURL(ctx, uri, response)
}

func (c *Client) CallURL(ctx context.Context, uri string, response any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return err
	}

	req.Header.Add("User-Agent", fmt.Sprintf("%s %s", c.agentName, c.agentEmail))
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

func (c *Client) GetCompanyTickers(ctx context.Context) ([]models.Company, error) {
	res := &models.CompanyTickersResponse{}
	var companies []models.Company
	err := c.CallURL(ctx, CompanyTickerCIKPath, res)
	if err != nil {
		return companies, err
	}
	for _, entry := range res.Data {
		if len(entry) < 4 {
			fmt.Println("insufficient data to unmarshall")
			continue
		}

		cik, ok := entry[0].(float64)
		if !ok {
			fmt.Println("invalik cik type")
			continue
		}

		name, ok := entry[1].(string)
		if !ok {
			fmt.Println("invalid name type")
			continue
		}

		ticker, ok := entry[2].(string)
		if !ok {
			fmt.Println("invalid ticker type")
			continue
		}

		var exch string
		if entry[3] != nil {
			exch, ok = entry[3].(string)
			if !ok {
				fmt.Println("invalid exch type")
				continue
			}
		}
		company := models.Company{
			CIK:      models.NumericCIK(cik),
			Name:     name,
			Ticker:   ticker,
			Exchange: exch,
		}
		if company.Exchange != "" {
			companies = append(companies, company)
		}
	}
	return companies, nil
}

func (c *Client) GetCompanyFacts(ctx context.Context, params *models.CompanyFactsParams) (*models.CompanyFactsResponse, error) {
	res := &models.CompanyFactsResponse{}
	err := c.Call(ctx, CompanyFactsPath, params, res)
	return res, err
}
