package polygon

import (
	"context"

	"github.com/jamesonhm/fingator/internal/polygon/models"
)

const (
	GroupedDailyPath = "/v2/aggs/grouped/locale/us/market/stocks/{date}"
)

func (c *Client) GroupedDailyBars(ctx context.Context, params *models.GroupedDailyParams) (*models.GroupedDailyResponse, error) {
	res := &models.GroupedDailyResponse{}
	err := c.Call(ctx, GroupedDailyPath, params, res)
	return res, err
}
