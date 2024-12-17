package models

import (
	"encoding/json"
	"fmt"
)

type NumericCIK int

func (c NumericCIK) Pad() string {
	return fmt.Sprintf("%010d", c)
}

type CompanyFactsParams struct {
	CIK NumericCIK `path:"cik_padded"`
}

type CompanyFactsResponse struct {
	CIK        int    `json:"cik"`
	EntityName string `json:"entityName"`
	Facts      struct {
		Data map[string]FactData `json:"us-gaap"`
	} `json:"facts"`
}

type FactData struct {
	Label string   `json:"label"`
	Units UnitData `json:"units"`
}

type UnitData struct {
	USD []UnitEntry `json:"USD"`
}

type UnitEntry struct {
	End          Date        `json:"end"`
	Value        json.Number `json:"val"`
	FiscalYear   int         `json:"fy"`
	FiscalPeriod string      `json:"fp"`
	Form         string      `json:"form"`
}
