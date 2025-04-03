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
			currStmt = &models.Statement{
				CIK:     cik,
				EndDate: line.EndD,
			}
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
		li := &models.LineItem{
			Tag:   line.Tag,
			Label: line.Label,
			Desc:  line.Description,
			Sheet: line.Statement,
			Units: line.Units,
			Value: val,
		}
		switch line.Category {
		case "NetIncome":
			currStmt.NetIncome = li
		case "DandA":
			currStmt.DandA = li
		case "Depreciation":
			currStmt.Depreciation = li
		case "Amortization":
			currStmt.Amortization = li
		case "NetCashOps":
			currStmt.NetCashOps = li
		case "CapEx":
			currStmt.CapEx = li
		case "DebtIssuance":
			currStmt.DebtIssuance = li
		case "DebtRepayment":
			currStmt.DebtRepayment = li
		case "Revenue":
			currStmt.Revenue = li
		case "EBIT":
			currStmt.EBIT = li
		case "TaxExpense":
			currStmt.TaxExpense = li
		case "PreTaxIncome":
			currStmt.PreTaxIncome = li
		case "EPS":
			currStmt.EPS = li
		case "Shares":
			currStmt.Shares = li
		case "CurrentAssets":
			currStmt.CurrentAssets = li
		case "CurrentLiabilities":
			currStmt.CurrentLiabilities = li
		case "ShareholderEquity":
			currStmt.ShareholderEquity = li
		}
	}
	return stmts, nil
}

func stmtInternals(stmts []*models.Statement) error {
	for _, stmt := range stmts {
		err := stmt.InternalCalcs(
			models.CalcTaxRate(),
			models.CalcNWC(),
		)
		if err != nil {
			return err
		}
	}
	return nil
}
