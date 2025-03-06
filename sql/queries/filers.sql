-- name: CreateFiler :one
INSERT INTO filers (cik, name)
VALUES (
    $1,
    $2
    )
ON CONFLICT ON CONSTRAINT filers_pkey DO UPDATE SET 
    name = EXCLUDED.name
RETURNING *;

-- name: GetFilers :many
SELECT cik, name
FROM filers;

-- name: CreateFiling :exec
INSERT INTO filings (
    accession_no,
    film_no,
    cik,
    filing_date
) VALUES (
    $1, $2, $3, $4
)
ON CONFLICT ON CONSTRAINT filings_pkey DO NOTHING;

-- name: CreateHolding :exec
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
);
