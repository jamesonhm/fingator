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
	return time.Now().Year() - f.LastFY()
}

func (f *FactData) LastFY() int {
	units, _ := f.LabelUnits()
	switch units {
	case "USD":
		return f.Units.USD[len(f.Units.USD)-1].FiscalYear
	case "PURE":
		return f.Units.Pure[len(f.Units.Pure)-1].FiscalYear
	case "SHARES":
		return f.Units.Shares[len(f.Units.Shares)-1].FiscalYear
	default:
		return 0
	}
}

func (f *FactData) LabelUnits() (string, error) {
	if len(f.Units.USD) > 0 {
		return "USD", nil
	} else if len(f.Units.Pure) > 0 {
		return "PURE", nil
	} else if len(f.Units.Shares) > 0 {
		return "SHARES", nil
	} else {
		return "", fmt.Errorf("Unknow units for label `%s`", f.Label)
	}
}

func (f *FactData) UnitEntries(lastn int) []UnitEntry {
	var l int
	units, _ := f.LabelUnits()
	switch units {
	case "USD":
		l = len(f.Units.USD)
		return f.Units.USD[l-min(l, lastn):]
	case "PURE":
		l = len(f.Units.Pure)
		return f.Units.Pure[l-min(l, lastn):]
	case "SHARES":
		l = len(f.Units.Shares)
		return f.Units.Shares[l-min(l, lastn):]
	default:
		return nil
	}
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
