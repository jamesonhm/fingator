package models

import "time"

//LatestFilingsPath = "https://www.sec.gov/cgi-bin/browse-edgar?action=getcurrent&CIK=&type={form_filter}&company=&dateb=&owner=include&start=0&count=100&output=atom"

type LatestFilingsParams struct {
	Action    Action     `query:"action"`
	CIK       *string    `query:"cik"`
	Type      *string    `query:"type"`
	Company   *string    `query:"company"`
	DateB     *string    `query:"dateb"`
	Ownership *Ownership `query:"owner"`
	Start     *int       `query:"start"`
	Count     *int       `query:"count"`
	Output    Output     `query:"output"`
}

type LatestFilingsResponse struct {
	Title   string        `xml:"title"`
	Updated time.Time     `xml:"updated"`
	Entries []FilingEntry `xml:"entry"`
}

type FilingEntry struct {
	Title   string    `xml:"title"`
	Link    Link      `xml:"link"`
	Summary string    `xml:"summary"`
	Updated time.Time `xml:"updated"`
	Form    Category  `xml:"category"`
	ID      string    `xml:"id"`
}

type Link struct {
	HRef string `xml:"href,attr"`
}

type Category struct {
	Type string `xml:"term,attr"`
}
