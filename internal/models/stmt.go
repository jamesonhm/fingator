package models

import (
	"fmt"
	"strconv"
	"time"
)

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
	CIK                     int32
	EndDate                 time.Time
	NetIncome               *LineItem
	DandA                   *LineItem
	Depreciation            *LineItem
	Amortization            *LineItem
	NetCashOps              *LineItem
	CapEx                   *LineItem
	DebtIssuance            *LineItem
	DebtRepayment           *LineItem
	Revenue                 *LineItem
	EBIT                    *LineItem
	TaxExpense              *LineItem
	PreTaxIncome            *LineItem
	EPS                     *LineItem
	Shares                  *LineItem
	TotalCurrentAssets      *LineItem
	CashEquivalents         *LineItem
	OtherNonOpAssets        *LineItem
	AccountsReceivable      *LineItem
	Inventory               *LineItem
	OtherOpAssets           *LineItem
	TotalCurrentLiabilities *LineItem
	AccountsPayable         *LineItem
	OtherOpLiabilities      *LineItem
	ShareholderEquity       *LineItem
	TaxRate                 *float64
	BalanceNWC              *float64
	DeltaNWC                *float64
	FCF                     *float64
	AnnualGrowth            *float64
}

func (s *Statement) String() string {
	var res string
	res += fmt.Sprintf("%d\n", s.CIK)
	res += fmt.Sprintf("%-12s\n", s.EndDate.Format("2006/01/02"))
	res += fmt.Sprintf("|%16s|%16.0f|\n", "NetIncome", s.NetIncome.Value.Float64())
	res += fmt.Sprintf("|%16s|%16.2f|\n", "TaxRate", floatNilCheck(s.TaxRate))
	res += fmt.Sprintf("|%16s|%16.0f|\n", "CurrAssets", s.TotalCurrentAssets.Value.Float64())
	res += fmt.Sprintf("|%16s|%16.0f|\n", "CashEquiv", s.CashEquivalents.Value.Float64())
	res += fmt.Sprintf("|%16s|%16.0f|\n", "OtherNonOpAssets", s.OtherNonOpAssets.Value.Float64())
	res += fmt.Sprintf("|%16s|%16.0f|\n", "AccountsPayable", s.AccountsPayable.Value.Float64())
	res += fmt.Sprintf("|%16s|%16.0f|\n", "OtherOpLiab.", s.OtherOpLiabilities.Value.Float64())
	res += fmt.Sprintf("|%16s|%16.0f|\n", "NWC", floatNilCheck(s.BalanceNWC))
	res += fmt.Sprintf("|%16s|%16.0f|\n", "ChangeNWC", floatNilCheck(s.DeltaNWC))
	return res
}

func floatNilCheck(f *float64) float64 {
	if f != nil {
		return *f
	}
	return 0
}

func (s *Statement) Calcs(opts ...CalcOption) error {
	var err error = nil
	for _, opt := range opts {
		s, err = opt(s)
		if err != nil {
			return err
		}
	}
	return nil
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

type CalcOption func(s *Statement) (*Statement, error)

func CalcTaxRate() CalcOption {
	return func(s *Statement) (*Statement, error) {
		if s.PreTaxIncome == nil || s.TaxExpense == nil {
			return s, fmt.Errorf("Missing Tax Rate Param - PreTax: %s, TaxExp: %s", s.PreTaxIncome.Value.StringValue(), s.TaxExpense.Value.StringValue())
		}
		tr := s.TaxExpense.Value.Float64() / s.PreTaxIncome.Value.Float64()
		s.TaxRate = &tr
		return s, nil
	}
}

func CalcBalanceNWC() CalcOption {
	return func(s *Statement) (*Statement, error) {
		if s.TotalCurrentAssets == nil || s.AccountsPayable == nil {
			return s, fmt.Errorf(
				"Missing NWC Param - Assets: %s, Liabilities: %s",
				s.TotalCurrentAssets.Value.StringValue(),
				s.AccountsPayable.Value.StringValue(),
			)
		}
		opassets := s.TotalCurrentAssets.Value.Float64()
		if s.CashEquivalents != nil {
			opassets -= s.CashEquivalents.Value.Float64()
		}
		if s.OtherNonOpAssets != nil {
			opassets -= s.OtherNonOpAssets.Value.Float64()
		}
		opliab := s.AccountsPayable.Value.Float64()
		if s.OtherOpLiabilities != nil {
			opliab += s.OtherOpLiabilities.Value.Float64()
		}
		nwc := opassets - opliab
		s.BalanceNWC = &nwc
		return s, nil
	}
}

func CalcDeltaNWC(prevStmt *Statement) CalcOption {
	return func(s *Statement) (*Statement, error) {
		if prevStmt.BalanceNWC == nil {
			return s, nil
		}
		if s.BalanceNWC == nil {
			return s, fmt.Errorf("Current Statement Missing NWC - %.2f", *s.BalanceNWC)
		}
		deltaNWC := *prevStmt.BalanceNWC - *s.BalanceNWC
		s.DeltaNWC = &deltaNWC
		return s, nil
	}
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
