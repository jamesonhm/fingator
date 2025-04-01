-- name: CreateFact :exec
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
ON CONFLICT ON CONSTRAINT facts_pkey DO NOTHING;

-- name: CompanyFacts :many
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
ORDER BY end_d, statement;
