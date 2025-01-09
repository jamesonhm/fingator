-- +goose Up
CREATE TABLE ohlc (
    ticker TEXT NOT NULL,
    ts TIMESTAMP NOT NULL,
    open NUMERIC NOT NULL,
    high NUMERIC NOT NULL,
    low NUMERIC NOT NULL,
    close NUMERIC NOT NULL,
    volume NUMERIC NOT NULL,
    CONSTRAINT ohlc_pkey PRIMARY KEY (ticker, ts)
);

-- +goose Down
DROP TABLE ohlc;
