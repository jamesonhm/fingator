-- name: GetUnmatchedCusips :many
SELECT DISTINCT(cusip)
FROM holdings
WHERE cusip NOT IN (
    SELECT cusip FROM cusips
);
