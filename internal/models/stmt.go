package models

import (
	"fmt"
	"strconv"
	"time"
)

type ValueHolder interface {
	StringValue() string
	Float64() float64
}

// USD represents dollar amount in cents
type USD int64

func USDFromStr(val string) USD {
	v, _ := strconv.ParseFloat(val, 64)
	return USD((v * 100) + 0.5)
}

func (u USD) StringValue() string {
	return fmt.Sprintf("$%.2f", float64(u)/100)
}

func (u USD) Float64() float64 {
	return float64(u) / 100
}

type Shares int64

func SharesFromStr(val string) Shares {
	v, _ := strconv.ParseInt(val, 10, 64)
	return Shares(v)
}

func (s Shares) StringValue() string {
	return fmt.Sprintf("%d", int64(s))
}

func (s Shares) Float64() float64 {
	return float64(s)
}

func NewStatement(cik int32, endd time.Time) *Statement {
	return &Statement{
		CIK:     cik,
		EndDate: endd,
		//Income:   make(map[string]LineItem),
		//Balance:  make(map[string]LineItem),
		//CashFlow: make(map[string]LineItem),
		//Calc:     StmtCalc{},
	}
}

type Statement struct {
	CIK                int32
	EndDate            time.Time
	NetIncome          *LineItem
	DandA              *LineItem
	Depreciation       *LineItem
	Amortization       *LineItem
	NetCashOps         *LineItem
	CapEx              *LineItem
	DebtIssuance       *LineItem
	DebtRepayment      *LineItem
	Revenue            *LineItem
	EBIT               *LineItem
	TaxExpense         *LineItem
	PreTaxIncome       *LineItem
	EPS                *LineItem
	Shares             *LineItem
	CurrentAssets      *LineItem
	CurrentLiabilities *LineItem
	ShareholderEquity  *LineItem
	//Income             map[string]LineItem
	//Balance            map[string]LineItem
	//CashFlow           map[string]LineItem
	TaxRate      *float64
	NWC          *float64
	FCF          *float64
	AnnualGrowth *float64
}

type LineItem struct {
	Tag   string
	Label string
	Desc  string
	Sheet string
	Units string
	Value ValueHolder
}

func (li LineItem) String() string {
	return fmt.Sprintf(
		"Tag: %s, Label: %s, Desc: %s...%s, Units: %s, Value: %v",
		li.Tag, li.Label, li.Desc[0:10], li.Desc[len(li.Desc)-11:len(li.Desc)], li.Units, li.Value,
	)
}

func (s *Statement) CalcTaxRate() {
	if s.PreTaxIncome == nil || s.TaxExpense == nil {

	}
	tr := s.TaxExpense.Value.Float64() / s.PreTaxIncome.Value.Float64()
	s.TaxRate = &tr
}

//// FCF = EBIT x (1- tax rate) + D&A + NWC – Capital expenditures
//func annualFCF(stmts []*models.Statement, logger *slog.Logger) {
//	var fcfs []float64
//	for i, stmt := range stmts {
//		fcf := stmt.Income["EBIT"]*(1-stmt.Calc.TaxRate) + stmt.CashFlow["DandA"] + stmt.Calc.NWC - stmt.CashFlow["CapEx"]
//	}
//}
//func annualGrowth(stmts []*models.Statement, logger *slog.Logger) {
//	var annualGrowth []float64
//	for i, stmt := range stmts {
//		if i == 0 {
//			continue
//		}
//
//	}
//}
