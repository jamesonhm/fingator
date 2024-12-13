package models

import (
	"time"
)

type Ticker struct {
	Active          bool      `json:"active"`
	Cik             string    `json:"cik"`
	CompositeFigi   string    `json:"composite_figi"`
	CurrencyName    string    `json:"currency_name"`
	LastUpdatedUtc  time.Time `json:"last_updated_utc"`
	Locale          string    `json:"locale"`
	Market          string    `json:"market"`
	Name            string    `json:"name"`
	PrimaryExchange string    `json:"primary_exchange"`
	ShareClassFigi  string    `json:"share_class_figi"`
	Ticker          string    `json:"ticker"`
	Type            string    `json:"type"`
}

type ListTickersParams struct {
	TickerEQ *string     `query:"ticker"`
	Type     *string     `query:"type"`
	Market   *AssetClass `query:"market"`
	Exchange *string     `query:"exchange"`
	CUSIP    *int        `query:"cusip"`
	CIK      *int        `query:"cik"`
	Date     *Date       `query:"date"`
	Active   *bool       `query:"active"`
	Search   *string     `query:"search"`
	Sort     *Sort       `query:"sort"`
	Order    *Order      `query:"order"`
	Limit    *int        `query:"limit"`
}

type ListTickersResponse struct {
	BaseResponse
	Results []Ticker `json:"results"`
}

type TickerDetailsParams struct {
	Ticker string `path:"ticker"`
	Date   *Date  `query:"date"`
}

type TickerDetailsResponse struct {
	BaseResponse
	Results Ticker `json:"results,omitempty"`
}
