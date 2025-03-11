-- name: GetUnmatchedCusips :many
SELECT DISTINCT(cusip)
FROM holdings
WHERE cusip NOT IN (
    SELECT cusip FROM cusips
);

-- name: AddCusip :exec
INSERT INTO cusips (
    cusip,
    security_name,
    ticker,
    exchange_code,
    security_type
) VALUES (
    $1, $2, $3, $4, $5
) ON CONFLICT ON CONSTRAINT cusips_pkey DO NOTHING;
