// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: facts.sql

package database

import (
	"context"
	"time"
)

const companyFacts = `-- name: CompanyFacts :many
SELECT cik,
    statement,
    category,
    tag,
    label,
    description,
    units,
    end_d,
    value,
    fiscal_year,
    fiscal_period,
    form
FROM facts
WHERE cik = $1
AND form = '10-K'
AND end_d > current_date - interval '5' year
ORDER BY end_d, statement
`

func (q *Queries) CompanyFacts(ctx context.Context, cik int32) ([]Fact, error) {
	rows, err := q.db.QueryContext(ctx, companyFacts, cik)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Fact
	for rows.Next() {
		var i Fact
		if err := rows.Scan(
			&i.Cik,
			&i.Statement,
			&i.Category,
			&i.Tag,
			&i.Label,
			&i.Description,
			&i.Units,
			&i.EndD,
			&i.Value,
			&i.FiscalYear,
			&i.FiscalPeriod,
			&i.Form,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const createFact = `-- name: CreateFact :exec
INSERT INTO facts (
    cik, 
    statement,
    category, 
    tag, 
    label, 
    description, 
    units, 
    end_d, 
    value, 
    fiscal_year, 
    fiscal_period, 
    form
)
VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
ON CONFLICT ON CONSTRAINT facts_pkey DO NOTHING
`

type CreateFactParams struct {
	Cik          int32
	Statement    string
	Category     string
	Tag          string
	Label        string
	Description  string
	Units        string
	EndD         time.Time
	Value        string
	FiscalYear   int32
	FiscalPeriod string
	Form         string
}

func (q *Queries) CreateFact(ctx context.Context, arg CreateFactParams) error {
	_, err := q.db.ExecContext(ctx, createFact,
		arg.Cik,
		arg.Statement,
		arg.Category,
		arg.Tag,
		arg.Label,
		arg.Description,
		arg.Units,
		arg.EndD,
		arg.Value,
		arg.FiscalYear,
		arg.FiscalPeriod,
		arg.Form,
	)
	return err
}
