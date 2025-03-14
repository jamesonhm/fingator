// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: comp.sql

package database

import (
	"context"
)

const createCompany = `-- name: CreateCompany :one
INSERT INTO companies (cik, name, ticker, exchange)
VALUES (
    $1,
    $2,
    $3,
    $4
    )
ON CONFLICT ON CONSTRAINT companies_pkey DO UPDATE SET 
    name = EXCLUDED.name,
    ticker = EXCLUDED.ticker,
    exchange = EXCLUDED.exchange
RETURNING cik, name, ticker, exchange
`

type CreateCompanyParams struct {
	Cik      int32
	Name     string
	Ticker   string
	Exchange string
}

func (q *Queries) CreateCompany(ctx context.Context, arg CreateCompanyParams) (Company, error) {
	row := q.db.QueryRowContext(ctx, createCompany,
		arg.Cik,
		arg.Name,
		arg.Ticker,
		arg.Exchange,
	)
	var i Company
	err := row.Scan(
		&i.Cik,
		&i.Name,
		&i.Ticker,
		&i.Exchange,
	)
	return i, err
}

const getExchangeCiks = `-- name: GetExchangeCiks :many
SELECT cik
FROM companies
WHERE exchange != ''
`

func (q *Queries) GetExchangeCiks(ctx context.Context) ([]int32, error) {
	rows, err := q.db.QueryContext(ctx, getExchangeCiks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int32
	for rows.Next() {
		var cik int32
		if err := rows.Scan(&cik); err != nil {
			return nil, err
		}
		items = append(items, cik)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
