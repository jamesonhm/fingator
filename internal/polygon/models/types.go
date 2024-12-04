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

type ListTickersResponse struct {
	BaseResponse
	Results []Ticker `json:"results"`
}
