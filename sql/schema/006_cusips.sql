-- +goose Up
CREATE TABLE cusips (
    cusip VARCHAR(255),
    security_name VARCHAR(255) NOT NULL,
    ticker VARCHAR(50) NOT NULL,
    exchange_code VARCHAR(10) NOT NULL,
    security_type VARCHAR(50),
    CONSTRAINT cusips_pkey PRIMARY KEY (cusip)
);

-- +goose Down
DROP TABLE cusips;
