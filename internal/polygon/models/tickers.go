package models

type Ticker struct {
	Active             bool           `json:"active"`
	Address            CompanyAddress `json:"address,omitempty"`
	Branding           Branding       `json:"branding,omitempty"`
	CIK                string         `json:"cik"`
	CompositeFigi      string         `json:"composite_figi"`
	CurrencyName       string         `json:"currency_name"`
	CurrencySymbol     string         `json:"currency_symbol,omitempty"`
	BaseCurrencyName   string         `json:"base_currency_name,omitempty"`
	BaseCurrencySymbol string         `json:"base_currency_symbol,omitempty"`
	Description        string         `json:"description,omitempty"`
	HomepageURL        string         `json:"homepage_url,omitempty"`
	ListDate           Date           `json:"list_date,omitempty"`
	MarketCap          float64        `json:"market_cap"`
	LastUpdatedUtc     Time           `json:"last_updated_utc"`
	Locale             string         `json:"locale"`
	Market             string         `json:"market"`
	Name               string         `json:"name"`
	PrimaryExchange    string         `json:"primary_exchange"`
	ShareClassFigi     string         `json:"share_class_figi"`
	Ticker             string         `json:"ticker"`
	Type               string         `json:"type"`
}

type CompanyAddress struct {
	Address1   string `json:"address1,omitempty"`
	Address2   string `json:"address2,omitempty"`
	City       string `json:"city,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
	State      string `json:"state,omitempty"`
}

type Branding struct {
	LogoURL string `json:"logo_url,omitempty"`
	IconURL string `json:"icon_url,omitempty"`
}

type ListTickersParams struct {
	TickerEQ *string     `query:"ticker"`
	Type     *string     `query:"type"`
	Market   *AssetClass `query:"market"`
	Exchange *string     `query:"exchange"`
	CUSIP    *int        `query:"cusip"`
	CIK      *int        `query:"cik"`
	Date     *Date       `query:"date"`
	Active   *bool       `query:"active"`
	Search   *string     `query:"search"`
	Sort     *Sort       `query:"sort"`
	Order    *Order      `query:"order"`
	Limit    *int        `query:"limit"`
}

type ListTickersResponse struct {
	BaseResponse
	Results []Ticker `json:"results"`
}

type TickerDetailsParams struct {
	Ticker string `path:"ticker"`
	Date   *Date  `query:"date"`
}

type TickerDetailsResponse struct {
	BaseResponse
	Results Ticker `json:"results,omitempty"`
}
