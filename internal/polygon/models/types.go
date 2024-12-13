package models

import (
	"encoding/json"
	"fmt"
	"strconv"
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

// PathFormat used to string format for use as a path parameter
func (d Date) PathFormat() string {
	return time.Time(d).Format("2006-01-02")
}

type Millis time.Time

// Unmarshaler interface to get timestamp string into Millis type
// https://pkg.go.dev/encoding/json#Unmarshaler
func (m *Millis) UnmarshalJSON(data []byte) error {
	d, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	*m = Millis(time.UnixMilli(d))
	return nil
}

func (m Millis) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(m).UnixMilli())
}

// Stringer Interface for print format
func (m Millis) String() string {
	return fmt.Sprintf("%v", time.Time(m).Format(time.DateTime))
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
