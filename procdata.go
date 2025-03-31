package main

import (
	"context"
	"log/slog"

	"github.com/jamesonhm/fingator/internal/database"
	"github.com/jamesonhm/fingator/internal/models"
)

type LineItem struct {
	Tag     string
	Label   string
	Desc    string
	Units   string
	Value   float64
	EndDate time.Time
}

type Statement struct {
	CIK        int32
	FiscalYear int
	Income     map[string]LineItem
	Balance    map[string]LineItem
	CashFlow   map[string]LineItem
}

func annualStatements(
	ctx context.Context,
	dbq *database.Queries,
	logger *slog.Logger,
	cik int32,
) ([]*models.Statement, error) {
	var stmts []*models.Statement
	var currStmt *models.Statement

	data, err := dbq.CompanyFacts(ctx, cik)
	if err != nil {
		return nil, err
	}

	for _, line := range data {
		if currStmt == nil || line.EndD != currStmt.EndDate {
			currStmt := models.NewStatement(cik, line.EndD)
			stmts = append(stmts, currStmt)
		}
		switch line.Statement {
		case "Income":
			currStmt.Income[line.Category] = models.LineItem{
				Tag:   line.Tag,
				Label: line.Label,
				Desc:  line.Description,
				Units: line.Units,
				Value: line.Value,
			}
		case "Balance":
			currStmt.Balance[line.Category] = models.LineItem{
				Tag:   line.Tag,
				Label: line.Label,
				Desc:  line.Description,
				Units: line.Units,
				Value: line.Value,
			}
		case "CashFlow":
			currStmt.CashFlow[line.Category] = models.LineItem{
				Tag:   line.Tag,
				Label: line.Label,
				Desc:  line.Description,
				Units: line.Units,
				Value: line.Value,
			}
		}
	}
	return stmts, nil
}
