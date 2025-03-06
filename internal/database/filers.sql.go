// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: filers.sql

package database

import (
	"context"
	"database/sql"
	"time"
)

const createFiler = `-- name: CreateFiler :one
INSERT INTO filers (cik, name)
VALUES (
    $1,
    $2
    )
ON CONFLICT ON CONSTRAINT filers_pkey DO UPDATE SET 
    name = EXCLUDED.name
RETURNING cik, name
`

type CreateFilerParams struct {
	Cik  int32
	Name string
}

func (q *Queries) CreateFiler(ctx context.Context, arg CreateFilerParams) (Filer, error) {
	row := q.db.QueryRowContext(ctx, createFiler, arg.Cik, arg.Name)
	var i Filer
	err := row.Scan(&i.Cik, &i.Name)
	return i, err
}

const createFiling = `-- name: CreateFiling :exec
INSERT INTO filings (
    accession_no,
    film_no,
    cik,
    filing_date
) VALUES (
    $1, $2, $3, $4
)
ON CONFLICT ON CONSTRAINT filings_pkey DO NOTHING
`

type CreateFilingParams struct {
	AccessionNo string
	FilmNo      int64
	Cik         int32
	FilingDate  time.Time
}

func (q *Queries) CreateFiling(ctx context.Context, arg CreateFilingParams) error {
	_, err := q.db.ExecContext(ctx, createFiling,
		arg.AccessionNo,
		arg.FilmNo,
		arg.Cik,
		arg.FilingDate,
	)
	return err
}

const createHolding = `-- name: CreateHolding :exec
INSERT INTO holdings (
    accession_no,
    name_of_issuer,
    class_title,
    cusip,
    value,
    shares,
    put_call
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
`

type CreateHoldingParams struct {
	AccessionNo  string
	NameOfIssuer string
	ClassTitle   string
	Cusip        string
	Value        int64
	Shares       int32
	PutCall      sql.NullString
}

func (q *Queries) CreateHolding(ctx context.Context, arg CreateHoldingParams) error {
	_, err := q.db.ExecContext(ctx, createHolding,
		arg.AccessionNo,
		arg.NameOfIssuer,
		arg.ClassTitle,
		arg.Cusip,
		arg.Value,
		arg.Shares,
		arg.PutCall,
	)
	return err
}

const getFilers = `-- name: GetFilers :many
SELECT cik, name
FROM filers
`

func (q *Queries) GetFilers(ctx context.Context) ([]Filer, error) {
	rows, err := q.db.QueryContext(ctx, getFilers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Filer
	for rows.Next() {
		var i Filer
		if err := rows.Scan(&i.Cik, &i.Name); err != nil {
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
