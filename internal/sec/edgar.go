package edgar

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jamesonhm/fingator/internal/encdec"
	"github.com/jamesonhm/fingator/internal/rate"
	"github.com/jamesonhm/fingator/internal/sec/models"
	"github.com/jamesonhm/fingator/internal/uri"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const (
	CompanyTickerCIKPath = "https://www.sec.gov/files/company_tickers_exchange.json"
	BaseDataURL          = "https://data.sec.gov"
	CompanyFactsPath     = "/api/xbrl/companyfacts/CIK{cik_padded}.json"
	LatestFilingsPath    = "https://www.sec.gov/cgi-bin/browse-edgar"
	MainURL              = "https://www.sec.gov"
)

type Client struct {
	agentName  string
	agentEmail string
	baseurl    string
	httpC      http.Client
	uriBuilder *uri.URIBuilder
	limiter    *rate.Limiter
}

// Create a new reference to an Edgar Client
// agentName and Email are sent as request headers to monitor requests
// period and count are used to rate limit requests to count/period
func New(agentName, agentEmail string, timeout time.Duration, period time.Duration, count int) *Client {
	return &Client{
		agentName:  agentName,
		agentEmail: agentEmail,
		baseurl:    BaseDataURL,
		httpC: http.Client{
			Timeout: timeout,
		},
		uriBuilder: uri.New(),
		limiter:    rate.New(period, count),
	}
}

// Call makes API call based on path and params
func (c *Client) Call(
	ctx context.Context,
	base string,
	path string,
	params,
	response any,
	decFunc models.DecFunc,
) error {
	uri := c.uriBuilder.EncodeParams(path, params)
	uri = base + uri
	fmt.Printf("client-call-uri: %s\n", uri)
	return c.CallURL(ctx, uri, response, decFunc)
}

// CallURL makes an API call based on a fully parameterized URL
func (c *Client) CallURL(ctx context.Context, uri string, response any, decFunc models.DecFunc) error {
	err := c.limiter.Wait(ctx)
	if err != nil {
		return err
	}
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

	if err := decFunc(resp, response); err != nil {
		return fmt.Errorf("error decoding: %w", err)
	}

	return nil
}

// CallURLopen makes an API call based on a fully parameterized URL and returns, not closes the resp body
func (c *Client) CallURLopen(ctx context.Context, uri string) (*http.Response, error) {
	err := c.limiter.Wait(ctx)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", fmt.Sprintf("%s %s", c.agentName, c.agentEmail))
	resp, err := c.httpC.Do(req)
	return resp, err
}

// Get a list of CIK's to tickers from SEC
func (c *Client) GetCompanyTickers(ctx context.Context) ([]models.Company, error) {
	res := &models.CompanyTickersResponse{}
	var companies []models.Company
	err := c.CallURL(ctx, CompanyTickerCIKPath, res, encdec.DecodeJsonResp)
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
		companies = append(companies, company)
	}
	return companies, nil
}

// Get a full list of company facts, each with their historical values from past submissions
func (c *Client) GetCompanyFacts(
	ctx context.Context,
	params *models.CompanyFactsParams,
) (*models.CompanyFactsResponse, error) {
	res := &models.CompanyFactsResponse{}
	err := c.Call(ctx, BaseDataURL, CompanyFactsPath, params, res, encdec.DecodeJsonResp)
	return res, err
}

func (c *Client) FetchFilings(
	ctx context.Context,
	params *models.BrowseEdgarParams,
) (*models.FetchFilingsResponse, error) {
	res := &models.FetchFilingsResponse{}
	err := c.Call(ctx, LatestFilingsPath, "", params, res, encdec.DecodeXmlResp)
	return res, err
}

func (c *Client) InfotableURLFromHTML(ctx context.Context, fe models.FilingEntry) (string, error) {
	url := fe.Link.Href.String()
	//fmt.Printf("URL: %s\n", url)
	res := &html.Node{}
	err := c.CallURL(ctx, url, res, encdec.DecodeHTMLResp)
	if err != nil {
		fmt.Printf("Error calling url")
		return "", err
	}

	for n := range res.Descendants() {
		if n.Type == html.TextNode && n.Parent.DataAtom == atom.A && strings.Contains(n.Data, ".xml") && !strings.Contains(n.Data, "primary") {
			//fmt.Printf("Data: %v, Parent Attr Href: %v\n", n.Data, n.Parent.Attr[0].Val)
			return MainURL + n.Parent.Attr[0].Val, nil
		}
	}
	return "", fmt.Errorf("link to xml filing not found")
}

func (c *Client) FetchHoldings(ctx context.Context, url string) (*models.FetchHoldingsResponse, error) {
	res := &models.FetchHoldingsResponse{}
	err := c.CallURL(ctx, url, res, encdec.DecodeXmlResp)
	return res, err
}

func (c *Client) File10kURLFromHTML(ctx context.Context, fe models.FilingEntry) (string, error) {
	url := fe.Link.Href.String()
	//fmt.Printf("URL: %s\n", url)
	res := &html.Node{}
	err := c.CallURL(ctx, url, res, encdec.DecodeHTMLResp)
	if err != nil {
		fmt.Printf("Error calling url")
		return "", err
	}
	fmt.Println("Filing HTML:", res)

	for n := range res.Descendants() {
		if n.Type == html.TextNode && n.Parent.DataAtom == atom.A && strings.Contains(n.Data, ".txt") && !strings.Contains(n.Data, "primary") {
			//fmt.Printf("Data: %v, Parent Attr Href: %v\n", n.Data, n.Parent.Attr[0].Val)
			return MainURL + n.Parent.Attr[0].Val, nil
		}
	}
	return "", fmt.Errorf("link to txt filing not found")
}

func (c *Client) Fetch10k(ctx context.Context, url string) (*http.Response, error) {
	resp, err := c.CallURLopen(ctx, url)
	return resp, err
}
