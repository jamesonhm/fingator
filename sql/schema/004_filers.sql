-- +goose Up
CREATE TABLE filers (
    cik INTEGER NOT NULL,
    name TEXT NOT NULL,
    CONSTRAINT filers_pkey PRIMARY KEY (cik)
);

INSERT INTO filers VALUES 
(1647251, "TCI Fund Management"),
(1072761, "Millenium Management"),
(1006438, "Appaloosa Management"),
(1603466, "Point72 Asset Management"),
(1135730, "Coatue Management"),
(1037389, "Renaissance Technologies"),
(1478735, "Two Sigma"),
(1103804, "Viking Global Investors");

-- +goose Down
DROP TABLE filers;
