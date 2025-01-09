package models

type GroupedDailyParams struct {
	Date       Date  `path:"date"`
	Adjusted   *bool `query:"adjusted"`
	IncludeOTC *bool `query:"include_otc"`
}

type GroupedDailyResponse struct {
	BaseResponse
	QueryCount  int  `json:"queryCount,omitempty"`
	ResultCount int  `json:"resultsCount,omitempty"`
	Adjusted    bool `json:"adjusted"`
	Results     []Agg
}

type Agg struct {
	Ticker       string  `json:"T"`
	Close        float64 `json:"c,omitempty"`
	High         float64 `json:"h,omitempty"`
	Low          float64 `json:"l,omitempty"`
	Transactions int32   `json:"n,omitempty"`
	Open         float64 `json:"o,omitempty"`
	Timestamp    Millis  `json:"t,omitempty"`
	Volume       float64 `json:"v,omitempty"`
	VWAP         float64 `json:"vw,omitempty"`
	OTC          bool    `json:"otc,omitempty"`
}
