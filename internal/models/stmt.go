package models

import "time"

type LineItem struct {
	Tag   string
	Label string
	Desc  string
	Units string
	Value float64
}

type Statement struct {
	CIK      int32
	EndDate  time.Time
	Income   map[string]LineItem
	Balance  map[string]LineItem
	CashFlow map[string]LineItem
}

func NewStatement(cik int32, endd time.Time) *Statement {
	return &Statement{
		CIK:      cik,
		EndDate:  endd,
		Income:   make(map[string]LineItem),
		Balance:  make(map[string]LineItem),
		CashFlow: make(map[string]LineItem),
	}
}
