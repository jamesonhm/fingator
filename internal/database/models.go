// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package database

import (
	"database/sql"
	"time"
)

type Company struct {
	Cik      int32
	Name     string
	Ticker   string
	Exchange string
}

type Cusip struct {
	Cusip        string
	SecurityName string
	Ticker       string
	ExchangeCode string
	SecurityType sql.NullString
}

type Fact struct {
	Cik          int32
	Category     string
	Tag          string
	Label        string
	Description  string
	Units        string
	EndD         time.Time
	Value        string
	FiscalYear   int32
	FiscalPeriod string
	Form         string
}

type Filer struct {
	Cik  int32
	Name string
}

type Filing struct {
	FilingID string
	Cik      int32
	Period   time.Time
}

type Holding struct {
	FilingID     string
	NameOfIssuer string
	Cusip        string
	Value        int64
	Shares       int32
}

type Ohlc struct {
	Ticker string
	Ts     time.Time
	Open   string
	High   string
	Low    string
	Close  string
	Volume string
}
