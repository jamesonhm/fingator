package models

import "encoding/json"

type IDType string

const (
	TypeCUSIP      IDType = "ID_CUSIP"
	TypeCUSIP8     IDType = "ID_CUSIP_8_CHR"
	TypeTicker     IDType = "TICKER"
	TypeBaseTicker IDType = "BASE_TICKER"
)

func (i IDType) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(i))
}
