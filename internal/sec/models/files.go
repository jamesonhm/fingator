package models

import (
	"encoding/json"
	"fmt"
)

//type FieldName string
//
//const (
//	CIK      FieldName = "cik"
//	Name     FieldName = "name"
//	Ticker   FieldName = "ticker"
//	Exchange FieldName = "exchange"
//)

type CompanyTickersResponse struct {
	Fields []string        `json:"fields"`
	Data   [][]interface{} `json:"data"`
}

type Company struct {
	CIK      NumericCIK
	Name     string
	Ticker   string
	Exchange string
}

func (c *Company) UnmarshallJSON(data []byte) error {
	var rawData []interface{}

	if err := json.Unmarshal(data, &rawData); err != nil {
		return err
	}

	if len(rawData) < 4 {
		return fmt.Errorf("insufficient data to unmarshall")
	}

	cik, ok := rawData[0].(float64)
	if !ok {
		return fmt.Errorf("invalik cik type")
	}
	c.CIK = NumericCIK(int64(cik))

	name, ok := rawData[1].(string)
	if !ok {
		return fmt.Errorf("invalid name type")
	}
	c.Name = name

	ticker, ok := rawData[2].(string)
	if !ok {
		return fmt.Errorf("invalid ticker type")
	}
	c.Ticker = ticker

	if rawData[3] != nil {
		exch, ok := rawData[3].(string)
		if !ok {
			return fmt.Errorf("invalid exch type")
		}
		c.Exchange = exch
	}
	return nil
}
