package main

import (
	//"database/sql"
	"context"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/jamesonhm/fingator/internal/database"
	"github.com/jamesonhm/fingator/internal/polygon"
	"github.com/jamesonhm/fingator/internal/polygon/models"
	edgar "github.com/jamesonhm/fingator/internal/sec"
	emodels "github.com/jamesonhm/fingator/internal/sec/models"
)

func runEdgarTickers(ctx context.Context, dbq *database.Queries, edgarClient edgar.Client, stdout, stderr io.Writer) {
	companies, err := edgarClient.GetCompanyTickers(ctx)
	if err != nil {
		fmt.Fprintf(stderr, "Error getting companies\n")
	}
	fmt.Fprintf(stdout, "no. of companies: %d\n", len(companies))

	for _, comp := range companies {
		_, err := dbq.CreateCompany(ctx, database.CreateCompanyParams{
			Cik:      int32(comp.CIK),
			Name:     comp.Name,
			Ticker:   comp.Ticker,
			Exchange: comp.Exchange,
		})
		if err != nil {
			fmt.Fprintf(stderr, "error adding company %+v: %v\n", comp, err)
		}
	}
}

func runEdgarFacts(ctx context.Context, dbq *database.Queries, edgarClient edgar.Client, stdout, stderr io.Writer) {
	ciks, err := dbq.GetExchangeCiks(ctx)
	if err != nil {
		fmt.Fprintf(stderr, "Error getting company ciks: %v\n", err)
		return
	}
	if len(ciks) == 0 {
		fmt.Fprintf(stderr, "Error, no CIK's found\n")
		return
	}

	for i, cik := range ciks {
		if i >= 5 {
			break
		}
		params := &emodels.CompanyFactsParams{
			CIK: emodels.NumericCIK(cik),
		}
		res, err := edgarClient.GetCompanyFacts(ctx, params)
		if err != nil {
			fmt.Fprintf(stderr, "Error getting company facts: %v\n", err)
		}

		dcf := &emodels.DCFData{}
		facts := edgar.FilterDCF(res, dcf)
		for _, fact := range facts {
			if len(fact.Units.USD) > 0 {
				fmt.Fprintf(
					stdout,
					"cik: %d, Category: %s, Tag: %s, Label: %s, units: %s, recent: %+v\n",
					cik, fact.Category, fact.Tag, fact.Label, "USD", fact.Units.USD[len(fact.Units.USD)-1],
				)
			} else if len(fact.Units.Pure) > 0 {
				fmt.Fprintf(
					stdout,
					"cik: %d, Category: %s, Tag: %s, Label: %s, units: %s, recent: %+v\n",
					cik, fact.Category, fact.Tag, fact.Label, "PURE", fact.Units.Pure[len(fact.Units.Pure)-1],
				)
			}
			fmt.Println()
		}
	}

	//fmt.Fprintf(stdout, "%+v\n", dcf)
}

func runEdgarFilings(ctx context.Context, getenv func(string) string, stdout, stderr io.Writer) {
	agentName := getenv("EDGAR_COMPANY_NAME")
	agentEmail := getenv("EDGAR_COMPANY_EMAIL")
	edgarClient := edgar.New(agentName, agentEmail, time.Second*10, 1)

	formType := "13F-HR"
	resCount := 100
	params := &emodels.BrowseEdgarParams{
		Action: emodels.GetCurrent,
		Type:   &formType,
		Count:  &resCount,
		Output: emodels.Atom,
	}
	res, err := edgarClient.FetchFilings(ctx, params)
	if err != nil {
		fmt.Fprintf(stderr, "error getting latest filings: %v\n", err)
	}

	for i, entry := range res.Entries {
		if i >= 3 {
			break
		}
		cik := string(entry.CIK())
		fmt.Fprintf(stdout, "Co (cik): %s (%s)\n", entry.Title, cik)
		path, _ := edgarClient.InfotableURLFromHTML(ctx, entry)
		fmt.Fprintf(stdout, "--%s\n\n", path)
	}

}

func runEdgarCompanyFilings(ctx context.Context, getenv func(string) string, stdout, stderr io.Writer) {
	agentName := getenv("EDGAR_COMPANY_NAME")
	agentEmail := getenv("EDGAR_COMPANY_EMAIL")
	edgarClient := edgar.New(agentName, agentEmail, time.Second*10, 1)

	formType := "13F-HR"
	resCount := 100
	cik := "0001471384"
	params := &emodels.BrowseEdgarParams{
		Action: emodels.GetCompany,
		Type:   &formType,
		Count:  &resCount,
		CIK:    &cik,
		Output: emodels.Atom,
	}
	compRes, err := edgarClient.FetchFilings(ctx, params)
	if err != nil {
		fmt.Fprintf(stderr, "error getting company filings: %v\n", err)
	}
	fmt.Fprintf(stdout, "%+v\n\n", compRes.CompanyInfo)
	for _, e := range compRes.Entries {
		fmt.Fprintf(stdout, "\nAccession: %v\n", e.AccessionNo())
		fmt.Fprintf(stdout, "Link: %s\n", e.Link.Href.String())
		path, _ := edgarClient.InfotableURLFromHTML(ctx, e)
		fmt.Fprintf(stdout, "--%s\n", path)
		holdings, err := edgarClient.FetchHoldings(ctx, path)
		if err != nil {
			fmt.Fprintf(stderr, "%v\n", err)
			continue
		}
		for i, h := range holdings.InfoTable {
			if i >= 10 {
				break
			}
			fmt.Fprintf(stdout, "*%+v\n", h)
		}
	}
}

func runPolyGrouped(ctx context.Context, dbq *database.Queries, polyClient polygon.Client, stdout, stderr io.Writer) {
	startEnd, err := dbq.OHLCStartEnd(ctx)
	if err != nil {
		fmt.Fprintf(stderr, "Error getting latest timestamp: %v\n", err)
	}
	fmt.Fprintf(stdout, "start: %v, end: %v\n", startEnd.Min, startEnd.Max)
	var minDate, maxDate *time.Time
	if start, ok := startEnd.Min.(time.Time); !ok {
		minDate = nil
	} else {
		minDate = &start
	}
	if end, ok := startEnd.Max.(time.Time); !ok {
		maxDate = nil
	} else {
		maxDate = &end
	}
	di := NewDateIter(5, minDate, maxDate, time.Now())
	for di.Next() {
		fmt.Fprintf(stdout, "next date: %v\n", di.Date)
		params := &models.GroupedDailyParams{
			Date: models.Date(di.Date),
		}
		res, err := polyClient.GroupedDailyBars(ctx, params)
		if err != nil {
			fmt.Fprintf(stderr, "Error happened here\n")
		}
		fmt.Fprintf(stdout, "result count: %d, status: %s\n", res.ResultCount, res.Status)
		for i, group := range res.Results {
			if i >= 5 {
				break
			}
			fmt.Fprintf(stdout, " * %+v\n", group)
			_, err := dbq.CreateTickerTimestamp(ctx, database.CreateTickerTimestampParams{
				Ticker: group.Ticker,
				Ts:     time.Time(group.Timestamp),
				Open:   strconv.FormatFloat(group.Open, 'f', 2, 64),
				High:   strconv.FormatFloat(group.High, 'f', 2, 64),
				Low:    strconv.FormatFloat(group.Low, 'f', 2, 64),
				Close:  strconv.FormatFloat(group.Close, 'f', 2, 64),
				Volume: strconv.FormatFloat(group.Volume, 'f', 2, 64),
			})
			if err != nil {
				fmt.Fprintf(stderr, "Error adding ticker/timestamp to db: %v\n", err)
			}
		}
	}
}
