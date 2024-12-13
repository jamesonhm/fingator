package polygon

import (
	"context"

	"github.com/jamesonhm/fingator/internal/polygon/iter"
	"github.com/jamesonhm/fingator/internal/polygon/models"
)

const (
	ListTickersPath   = "/v3/reference/tickers"
	TickerDetailsPath = "/v3/reference/tickers/{ticker}"
)

// ListTickers retrieves a list of tickers
// This method returns an iterator for accessing the results sequentially, across multiple pages if required
//
// iter := c.ListTickers(context.TODO(), params)
//
//	for iter.Next() {
//		log.Print(iter.Item()) // do something with each item
//	}
//	if iter.Err() != nil {
//		return iter.Err()
//	}
func (c *Client) ListTickers(ctx context.Context, params *models.ListTickersParams) *iter.Iter[models.Ticker] {
	return iter.NewIter(ctx, ListTickersPath, params, func(uri string) (iter.ListResponse, []models.Ticker, error) {
		res := &models.ListTickersResponse{}
		err := c.CallURL(ctx, uri, res)
		return res, res.Results, err
	})
}

func (c *Client) GetTickerDetails(ctx context.Context, params *models.TickerDetailsParams) (*models.TickerDetailsResponse, error) {
	res := &models.TickerDetailsResponse{}
	err := c.Call(ctx, TickerDetailsPath, params, res)
	return res, err
}
