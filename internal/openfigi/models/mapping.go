package models

type MappingRequest struct {
	IDType   IDType       `json:"idType"`
	IDValue  string       `json:"idValue"`
	ExchCode ExchangeCode `json:"exchCode"`
}

type MappingResponse struct {
	Data    []Object `json:"data"`
	Error   string   `json:"error,omitempty"`
	Warning string   `json:"warning,omitempty"`
}

type Object struct {
	FIGI         string `json:"figi"`
	SecurityType string `json:"securityType,omitempty"`
	MarketSector string `json:"marketSector,omitempty"`
	Ticker       string `json:"ticker,omitempty"`
	Name         string `json:"name,omitempty"`
	ExchangeCode string `json:"exchCode,omitempty"`
	MetaData     string `json:"metadata,omitempty"`
}
