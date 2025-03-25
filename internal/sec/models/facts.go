package models

import (
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"time"
)

type CompanyFactsParams struct {
	CIK NumericCIK `path:"cik_padded"`
}

type CIKWrapper int

// TODO: This could be flattened with a custom Unmarshal, eliminating all the switching on units...
// https://go.dev/play/p/83VHShfE5rI
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
	findyr := func(ue []UnitEntry) int {
		maxyr := 0
		for _, e := range ue {
			if e.FiscalYear > maxyr {
				maxyr = e.FiscalYear
			}
		}
		return maxyr
	}
	units, _ := f.LabelUnits()
	switch units {
	case "USD":
		return findyr(f.Units.USD)
		//return f.Units.USD[len(f.Units.USD)-1].FiscalYear
	case "PURE":
		return findyr(f.Units.Pure)
		//return f.Units.Pure[len(f.Units.Pure)-1].FiscalYear
	case "SHARES":
		return findyr(f.Units.Shares)
		//return f.Units.Shares[len(f.Units.Shares)-1].FiscalYear
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

func (f *FactData) Filter() error {
	var ue []UnitEntry
	minYr := time.Now().Year() - 10
	units, err := f.LabelUnits()
	if err != nil {
		return err
	}
	delf := func(v UnitEntry) bool {
		return (v.Form != "10-K" && v.Form != "10-Q" && v.Form != "10-Q/A") || v.FiscalYear < minYr
	}
	switch units {
	case "USD":
		ue = f.Units.USD
		ue = slices.DeleteFunc(ue, delf)
		if len(ue) == 0 {
			return fmt.Errorf("All elements filtered out")
		}
		f.Units.USD = ue
	case "PURE":
		ue = f.Units.Pure
		ue = slices.DeleteFunc(ue, delf)
		if len(ue) == 0 {
			return fmt.Errorf("All elements filtered out")
		}
		f.Units.Pure = ue
	case "SHARES":
		ue = f.Units.Shares
		ue = slices.DeleteFunc(ue, delf)
		if len(ue) == 0 {
			return fmt.Errorf("All elements filtered out")
		}
		f.Units.Shares = ue
	}
	return nil
}

func (f *FactData) UnitEntries() []UnitEntry {
	units, _ := f.LabelUnits()
	switch units {
	case "USD":
		return f.Units.USD
	case "PURE":
		return f.Units.Pure
	case "SHARES":
		return f.Units.Shares
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
	Sheet    string
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
