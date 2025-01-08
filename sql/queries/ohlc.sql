-- name: CreateTickerTimestamp :one
INSERT INTO ohlc (ticker, ts, open, high, low, close, num_trans, volume, vol_weighted)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $9
    )
ON CONFLICT ON CONSTRAINT ohlc_pkey DO UPDATE SET (
    open = EXCLUDED.open,
    high = EXCLUDED.high,
    low = EXCLUDED.low,
    close = EXCLUDED.close,
    num_trans = EXCLUDED.num_trans,
    volume = EXCLUDED.volume,
    vol_weighted = EXCLUDED.vol_weighted
)
RETURNING *;
