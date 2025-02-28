-- +goose Up
CREATE TABLE filings (
    filing_id VARCHAR(255) PRIMARY KEY,
    cik INTEGER NOT NULL REFERENCES filers(cik),
    period DATE NOT NULL,
);

-- +goose Down
DROP TABLE filings;
