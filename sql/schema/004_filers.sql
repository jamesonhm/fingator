-- +goose Up
CREATE TABLE filers (
    cik INTEGER NOT NULL,
    category TEXT NOT NULL,
    CONSTRAINT filers_pkey PRIMARY KEY (cik)
);

-- +goose Down
DROP TABLE filers;
