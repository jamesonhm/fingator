package main

import (
	"context"
	"log/slog"

	"github.com/jamesonhm/fingator/internal/database"
	"github.com/jamesonhm/fingator/internal/models"
)

func annualStatements(
	ctx context.Context,
	dbq *database.Queries,
	logger *slog.Logger,
	cik int32,
) ([]*models.Statement, error) {
	var stmts []*models.Statement
	var currStmt *models.Statement
	var val models.ValueHolder

	data, err := dbq.CompanyFacts(ctx, cik)
	if err != nil {
		return nil, err
	}

	for _, line := range data {
		if currStmt == nil || line.EndD != currStmt.EndDate {
			currStmt = models.NewStatement(cik, line.EndD)
			stmts = append(stmts, currStmt)
		}
		switch line.Units {
		case "USD":
			val = models.USDFromStr(line.Value)
		case "USD/shares":
			val = models.USDFromStr(line.Value)
		case "SHARES":
			val = models.SharesFromStr(line.Value)
		}
		li := models.LineItem{
			Tag:   line.Tag,
			Label: line.Label,
			Desc:  line.Description,
			Units: line.Units,
			Value: val,
		}
		key := line.Category
		switch line.Statement {
		case "Income":
			currStmt.Income[key] = li
		case "Balance":
			currStmt.Balance[key] = li
		case "CashFlow":
			currStmt.CashFlow[key] = li
		}
	}
	return stmts, nil
}

func calcTaxRate(ctx context.Context, stmts []*models.Statement, logger *slog.Logger) {
	for _, stmt := range stmts {
		pretax, ok := stmt.Income["PreTaxIncome"]
		taxrate, ok := stmt.Income["TaxExpense"]
		if !ok {
			logger.LogAttrs(
				ctx,
				slog.LevelError,
				"Missing Tax Rate Param",
				slog.Int("CIK", int(stmt.CIK)),
				slog.Time("Date", stmt.EndDate),
				slog.Float64("PreTaxIncome", pretax.Value.Float64()),
				slog.Float64("TaxExpense", taxrate.Value.Float64()),
			)
			continue
		}
		tr := taxrate.Value.Float64() / pretax.Value.Float64()
		stmt.Calc.TaxRate = &tr
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
