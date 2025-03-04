-- +goose Up
CREATE TABLE filings (
    filing_id VARCHAR(255),
    cik INTEGER NOT NULL REFERENCES filers(cik),
    period DATE NOT NULL,
    CONSTRAINT filings_pkey PRIMARY KEY (filing_id)
);

-- +goose Down
DROP TABLE filings;
