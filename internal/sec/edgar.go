package edgar

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jamesonhm/fingator/internal/encdec"
	"github.com/jamesonhm/fingator/internal/sec/models"
	"github.com/jamesonhm/fingator/internal/uri"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"golang.org/x/time/rate"
)

const (
	CompanyTickerCIKPath = "https://www.sec.gov/files/company_tickers_exchange.json"
	BaseURL              = "https://data.sec.gov"
	CompanyFactsPath     = "/api/xbrl/companyfacts/CIK{cik_padded}.json"
	LatestFilingsPath    = "https://www.sec.gov/cgi-bin/browse-edgar"
)

type Client struct {
	agentName  string
	agentEmail string
	baseurl    string
	httpC      http.Client
	uriBuilder *uri.URIBuilder
	limiter    *rate.Limiter
}

func New(agentName, agentEmail string, timeout time.Duration, reqsPerSec rate.Limit) Client {
	return Client{
		agentName:  agentName,
		agentEmail: agentEmail,
		baseurl:    BaseURL,
		httpC: http.Client{
			Timeout: timeout,
		},
		uriBuilder: uri.New(),
		limiter:    rate.NewLimiter(reqsPerSec, 1),
	}
}

// Call makes API call based on path and params
func (c *Client) Call(ctx context.Context, base string, path string, params, response any, decFunc models.DecFunc) error {
	uri := c.uriBuilder.EncodeParams(path, params)
	uri = base + uri
	//fmt.Printf("client-call-uri: %s\n", uri)
	return c.CallURL(ctx, uri, response, decFunc)
}

func (c *Client) CallURL(ctx context.Context, uri string, response any, decFunc models.DecFunc) error {
	c.limiter.Wait(ctx)
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
	//if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
	if err := decFunc(resp, response); err != nil {
		return fmt.Errorf("error decoding json: %w", err)
	}

	return nil
}

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
		if company.Exchange != "" {
			companies = append(companies, company)
		}
	}
	return companies, nil
}

func (c *Client) GetCompanyFacts(ctx context.Context, params *models.CompanyFactsParams) (*models.CompanyFactsResponse, error) {
	res := &models.CompanyFactsResponse{}
	err := c.Call(ctx, BaseURL, CompanyFactsPath, params, res, encdec.DecodeJsonResp)
	return res, err
}

func (c *Client) FetchLatestFiling(ctx context.Context, params *models.LatestFilingsParams) (*models.LatestFilingsResponse, error) {
	res := &models.LatestFilingsResponse{}
	err := c.Call(ctx, LatestFilingsPath, "", params, res, encdec.DecodeXmlResp)
	return res, err
}

func (c *Client) InfotableURLFromHTML(ctx context.Context, fe models.FilingEntry) (string, error) {
	url := fe.Link.HRef.String()
	fmt.Printf("URL: %s\n", url)
	res := &html.Node{}
	err := c.CallURL(ctx, url, res, encdec.DecodeHTMLResp)
	if err != nil {
		fmt.Printf("Error calling url")
		return "", err
	}

	fmt.Printf("NodeRes: %+v\n", res)

	for n := range res.Descendants() {
		//fmt.Printf(n.Data)
		//fmt.Printf("Atom: %v\n", n.DataAtom)
		if n.Type == html.TextNode && n.Parent.DataAtom == atom.A && strings.Contains(n.Data, ".xml") && !strings.Contains(n.Data, "primary") {
			fmt.Printf("Data: %v, Parent Attr Href: %v\n", n.Data, n.Parent.Attr[0].Val)
		}
	}
	//pathParts[len(pathParts)-1] = "infotable.xml"
	//infotablePath := strings.Join(pathParts, "/")
	//u := &url.URL{
	//	Scheme: fe.Link.HRef.Scheme,
	//	Host:   fe.Link.HRef.Host,
	//	Path:   infotablePath,
	//}
	return "", nil
}
