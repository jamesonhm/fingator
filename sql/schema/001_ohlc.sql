-- +goose Up
CREATE TABLE ohlc (
    ticker TEXT NOT NULL,
    ts TIMESTAMP NOT NULL,
    open MONEY,
    high MONEY,
    low MONEY,
    close MONEY,
    num_trans INTEGER,
    volume NUMERIC,
    vol_weighted NUMERIC,
    CONSTRAINT ohlc_pkey PRIMARY KEY (ticker, ts)
);

-- +goose Down
DROP TABLE ohlc;
