package models

import (
	"time"
)

type AssetClass string

const (
	AssetStocks  AssetClass = "stocks"
	AssetOptions AssetClass = "options"
	AssetCrypto  AssetClass = "crypto"
	AssetFx      AssetClass = "fx"
	AssetOTC     AssetClass = "otc"
	AssetIndices AssetClass = "indices"
)

type Date time.Time

func (d Date) Format() string {
	return time.Time(d).Format("2006-01-02")
}

type Order string

const (
	Asc  Order = "asc"
	Desc Order = "desc"
)

// Sort is a query param type
type Sort string

const (
	TickerSymbol    Sort = "ticker"
	Name            Sort = "name"
	Market          Sort = "market"
	Locale          Sort = "locale"
	PrimaryExchange Sort = "primary_exchange"
	Type            Sort = "type"
	CurrencySymbol  Sort = "currency_symbol"
	CurrencyName    Sort = "currency_name"
	Timestamp       Sort = "timestamp"
)
