-- +goose Up
CREATE TABLE companies (
    cik INTEGER NOT NULL,
    name TEXT NOT NULL,
    ticker TEXT NOT NULL,
    exchange TEXT NOT NULL,
    CONSTRAINT companies_pkey PRIMARY KEY (cik)
);

-- +goose Down
DROP TABLE companies;
