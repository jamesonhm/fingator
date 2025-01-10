-- name: CreateTickerTimestamp :one
INSERT INTO ohlc (ticker, ts, open, high, low, close, volume)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
    )
ON CONFLICT ON CONSTRAINT ohlc_pkey DO UPDATE SET 
    open = EXCLUDED.open,
    high = EXCLUDED.high,
    low = EXCLUDED.low,
    close = EXCLUDED.close,
    volume = EXCLUDED.volume
RETURNING *;

-- name: OHLCStartEnd :one
SELECT MIN(ts), MAX(ts) FROM ohlc;

