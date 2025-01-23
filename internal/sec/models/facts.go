package models

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type CompanyFactsParams struct {
	CIK NumericCIK `path:"cik_padded"`
}

type CIKWrapper int

type CompanyFactsResponse struct {
	CIK        CIKWrapper `json:"cik"`
	EntityName string     `json:"entityName"`
	Facts      struct {
		USGAAP map[string]FactData `json:"us-gaap"`
		DEI    map[string]FactData `json:"dei"`
	} `json:"facts"`
}

type FactData struct {
	Label       string   `json:"label"`
	Description string   `json:"description"`
	Units       UnitData `json:"units"`
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

func (c *CIKWrapper) UnmarshalJSON(b []byte) error {
	var i interface{}
	if err := json.Unmarshal(b, &i); err != nil {
		return err
	}
	switch t := i.(type) {
	case int:
		*c = CIKWrapper(t)
	case float64:
		*c = CIKWrapper(int(t))
	case string:
		if s, err := strconv.Atoi(t); err == nil {
			*c = CIKWrapper(s)
		} else {
			return err
		}
	default:
		return fmt.Errorf("cik response type not handled: %T", t)
	}
	return nil
}
