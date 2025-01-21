-- name: CreateFact :exec
INSERT INTO facts (cik, category, tag, label, description, units, end_d, value, fiscal_year, fiscal_period, form)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10,
    $11
    )
ON CONFLICT ON CONSTRAINT facts_pkey DO NOTHING;
