-- name: CreateCompany :one
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
RETURNING *;

-- name: GetExchangeCiks :many
SELECT cik
FROM companies
WHERE exchange != '';
