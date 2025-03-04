-- +goose Up
CREATE TABLE holdings (
    filing_id VARCHAR(255) NOT NULL REFERENCES filings(filing_id),
    name_of_issuer VARCHAR(255) NOT NULL,
    cusip VARCHAR(12) NOT NULL,
    value BIGINT NOT NULL,
    shares INTEGER NOT NULL
);

-- +goose Down
DROP TABLE holdings;
