-- +goose Up
CREATE TABLE filings (
    accession_no VARCHAR(255) NOT NULL,
    film_no BIGINT NOT NULL,
    cik INTEGER NOT NULL REFERENCES filers(cik),
    filing_date DATE NOT NULL,
    CONSTRAINT filings_pkey PRIMARY KEY (accession_no)
);

-- +goose Down
DROP TABLE filings;
