-- name: CreateCompnay :one
INSERT INTO companies (cik, name, ticker, exchange)
VALUES (
    $1,
    $2,
    $3,
    $4
    )
ON CONFLICT  DO UPDATE SET 

RETURNING *;

