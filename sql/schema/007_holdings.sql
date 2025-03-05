-- +goose Up
CREATE TABLE holdings (
    accession_no VARCHAR(255) NOT NULL REFERENCES filings(accession_no),
    name_of_issuer VARCHAR(255) NOT NULL,
    cusip VARCHAR(12) NOT NULL,
    value BIGINT NOT NULL,
    shares INTEGER NOT NULL
);

-- +goose Down
DROP TABLE holdings;
