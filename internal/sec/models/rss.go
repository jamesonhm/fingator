package models

import (
	"encoding/xml"
	"net/url"
	"strings"
	"time"
)

//LatestFilingsPath = "https://www.sec.gov/cgi-bin/browse-edgar?action=getcurrent&CIK=&type={form_filter}&company=&dateb=&owner=include&start=0&count=100&output=atom"

type BrowseEdgarParams struct {
	Action    Action     `query:"action"`
	FileNum   *string    `query:"filenum"`
	CIK       *string    `query:"CIK"`
	Type      *string    `query:"type"`
	Company   *string    `query:"company"`
	DateB     *string    `query:"dateb"`
	Ownership *Ownership `query:"owner"`
	Start     *int       `query:"start"`
	Count     *int       `query:"count"`
	Output    Output     `query:"output"`
}

type FetchFilingsResponse struct {
	CompanyInfo CompanyInfo   `xml:"company-info,omitempty"`
	Entries     []FilingEntry `xml:"entry"`
	Title       string        `xml:"title"`
	Updated     time.Time     `xml:"updated"`
}

type CompanyInfo struct {
	Addresses struct {
		Address []struct {
			Type    string `xml:"type,attr"`
			City    string `xml:"city"`
			State   string `xml:"state"`
			Street1 string `xml:"street1"`
			Street2 string `xml:"street2"`
			Zip     string `xml:"zip"`
			Phone   string `xml:"phone"`
		} `xml:"address"`
	} `xml:"addresses"`
	CIK                string `xml:"cik"`
	ConformedName      string `xml:"conformed-name"`
	StateLocation      string `xml:"state-location"`
	StateIncorporation string `xml:"state-of-incorporation"`
}

type FilingEntry struct {
	Title   string    `xml:"title"`
	Link    Link      `xml:"link"`
	Summary string    `xml:"summary"`
	Updated time.Time `xml:"updated"`
	Form    Category  `xml:"category"`
	ID      string    `xml:"id"`
	Content Content   `xml:"content,omitempty"`
}

func (fe *FilingEntry) AccessionNo() string {
	return strings.Split(fe.ID, "=")[1]
}

func (fe *FilingEntry) CIK() string {
	return strings.Split(fe.AccessionNo(), "-")[0]
}

type Content struct {
	AccessionNumber string `xml:"accession-number,omitempty"`
	FileNumber      string `xml:"file-number,omitempty"`
	FilingDate      string `xml:"filing-date,omitempty"`
	FilingHref      string `xml:"filing-href,omitempty"`
	FilingType      string `xml:"filing-type,omitempty"`
	FilmNumber      string `xml:"film-number,omitempty"`
	Amend           string `xml:"amend,omitempty"`
}

type Link struct {
	Href url.URL `xml:"href,attr"`
}

type Category struct {
	Type string `xml:"term,attr"`
}

func (l *Link) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "href":
			parsedURL, err := url.Parse(attr.Value)
			if err != nil {
				return err
			}
			l.Href = *parsedURL
		}
	}
	d.Skip()
	return nil
}

type InformationTable struct {
	InfoTable []Holding `xml:"infoTable"`
}

type Holding struct {
	NameOfIssuer   string `xml:"nameOfIssuer"`
	TitleOfClass   string `xml:"titleOfClass"`
	CUSIP          string `xml:"cusip"`
	Value          int    `xml:"value"`
	SharesOrPrnAmt struct {
		Amount int    `xml:"sshPrnamt"`
		Type   string `xml:"sshPrnamtType"`
	}
}
