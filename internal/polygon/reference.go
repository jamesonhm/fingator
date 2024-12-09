package polygon

import (
	"context"
	"fmt"

	"github.com/jamesonhm/fingator/internal/polygon/models"
)

const (
	ListTickersPath = "/v3/reference/tickers"
)

// ListTickers retrieves a list of tickers
// This methon returns an iterator for accessing the results sequentially, across multiple pages if required
func (c *Client) ListTickers(ctx context.Context, params *models.ListTickersParams) error {
	res := &models.ListTickersResponse{}
	err := c.Call(ctx, ListTickersPath, params, res)
	fmt.Printf("List Tickers Response: \n")
	for _, r := range res.Results {
		fmt.Printf("%v\n", r)
	}
	return err
}
