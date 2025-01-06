package models

import (
	"encoding/json"
)

type CompanyFactsParams struct {
	CIK NumericCIK `path:"cik_padded"`
}

type CompanyFactsResponse struct {
	CIK        int    `json:"cik"`
	EntityName string `json:"entityName"`
	Facts      struct {
		USGAAP map[string]FactData `json:"us-gaap"`
		DEI    map[string]FactData `json:"dei"`
	} `json:"facts"`
}

type FactData struct {
	Label string   `json:"label"`
	Units UnitData `json:"units"`
}

type UnitData struct {
	USD    []UnitEntry `json:"USD,omitempty"`
	Pure   []UnitEntry `json:"pure,omitempty"`
	Shares []UnitEntry `json:"shares,omitempty"`
}

type UnitEntry struct {
	End          Date        `json:"end"`
	Value        json.Number `json:"val"`
	FiscalYear   int         `json:"fy"`
	FiscalPeriod string      `json:"fp"`
	Form         string      `json:"form"`
}

type FilteredFact struct {
	Category string
	Tag      string
	FactData
}

type DCFData struct {
	CashFlow         FactData
	CapEx            FactData
	Revenue          FactData
	NetIncome        FactData
	OperatingExpense FactData
}
