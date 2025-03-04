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

INSERT INTO filings (
    filing_id,
    cik,
    period
) VALUES (
    $1, $2, $3
)
ON CONFLICT ON CONSTRAINT filings_pkey DO NOTHING;
