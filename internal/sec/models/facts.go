package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
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

func (f *FactData) Age() int {
	var fy int
	switch f.UnitLabel() {
	case "USD":
		fy = f.Units.USD[len(f.Units.USD)-1].FiscalYear
	case "PURE":
		fy = f.Units.Pure[len(f.Units.Pure)-1].FiscalYear
	case "SHARES":
		fy = f.Units.Shares[len(f.Units.Shares)-1].FiscalYear
	}
	return time.Now().Year() - fy
}

func (f *FactData) UnitLabel() string {
	if len(f.Units.USD) > 0 {
		return "USD"
	} else if len(f.Units.Pure) > 0 {
		return "PURE"
	} else {
		return "SHARES"
	}
}

func (f *FactData) UnitEntries() []UnitEntry {

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
