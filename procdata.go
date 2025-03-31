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
			currStmt := models.NewStatement(cik, line.EndD)
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
		switch line.Statement {
		case "Income":
			currStmt.Income[line.Category] = li
		case "Balance":
			currStmt.Balance[line.Category] = li
		case "CashFlow":
			currStmt.CashFlow[line.Category] = li
		}
	}
	return stmts, nil
}
