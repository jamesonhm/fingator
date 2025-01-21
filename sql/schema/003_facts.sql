-- +goose Up
CREATE TABLE facts (
    cik INTEGER NOT NULL,
    category TEXT NOT NULL,
    tag TEXT NOT NULL,
    label TEXT NOT NULL,
    description TEXT NOT NULL,
    units TEXT NOT NULL,
    end_d DATE NOT NULL,
    value NUMERIC NOT NULL,
    fiscal_year INTEGER NOT NULL,
    fiscal_period TEXT NOT NULL,
    form TEXT NOT NULL,
    CONSTRAINT facts_pkey PRIMARY KEY (cik, category, end_d)
);

-- +goose Down
DROP TABLE facts;
