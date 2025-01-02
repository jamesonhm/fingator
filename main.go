package main

import (
	//"database/sql"
	"context"
	"fmt"
	"io"
	"time"

	//"log"
	"os"

	"github.com/jamesonhm/fingator/internal/polygon"
	"github.com/jamesonhm/fingator/internal/polygon/models"
	edgar "github.com/jamesonhm/fingator/internal/sec"
	emodels "github.com/jamesonhm/fingator/internal/sec/models"
	"github.com/joho/godotenv"
)

func runEdgarTickers(ctx context.Context, getenv func(string) string, stdout, stderr io.Writer) {
	agentName := getenv("EDGAR_COMPANY_NAME")
	agentEmail := getenv("EDGAR_COMPANY_EMAIL")
	edgarClient := edgar.New(agentName, agentEmail, time.Second*10)

	companies, err := edgarClient.GetCompanyTickers(ctx)
	if err != nil {
		fmt.Fprintf(stderr, "Error getting companies\n")
	}

	fmt.Fprintf(stdout, "no. of companies: %d\n", len(companies))
}

func runEdgarFacts(ctx context.Context, getenv func(string) string, stdout, stderr io.Writer) {
	agentName := getenv("EDGAR_COMPANY_NAME")
	agentEmail := getenv("EDGAR_COMPANY_EMAIL")
	edgarClient := edgar.New(agentName, agentEmail, time.Second*10)

	comps, err := edgarClient.GetCompanyTickers(ctx)
	if err != nil {
		fmt.Fprintf(stderr, "Error getting company tickers\n")
	}

	for i, comp := range comps {
		if i >= 5 {
			break
		}

		fmt.Fprintf(stdout, "comp: %s, ticker: %s, CIK: %d\n", comp.Name, comp.Ticker, comp.CIK)
		params := &emodels.CompanyFactsParams{
			CIK: comp.CIK,
		}
		res, err := edgarClient.GetCompanyFacts(ctx, params)
		if err != nil {
			fmt.Fprintf(stderr, "Error getting company facts\n")
		}

		dcf := &emodels.DCFData{}
		facts := edgar.FilterDCF(res, dcf)
		for _, fact := range facts {
			if len(fact.Units.USD) > 0 {
				fmt.Printf("Category: %s, Tag: %s, Label: %s, recent: %+v\n", fact.Category, fact.Tag, fact.Label, fact.Units.USD[len(fact.Units.USD)-1])
			} else if len(fact.Units.Pure) > 0 {
				fmt.Printf("Category: %s, Tag: %s, Label: %s, recent: %+v\n", fact.Category, fact.Tag, fact.Label, fact.Units.Pure[len(fact.Units.Pure)-1])
			}
			fmt.Println()
		}
	}

	//fmt.Fprintf(stdout, "%+v\n", dcf)
}

func runEdgarFilings(ctx context.Context, getenv func(string) string, stdout, stderr io.Writer) {
	agentName := getenv("EDGAR_COMPANY_NAME")
	agentEmail := getenv("EDGAR_COMPANY_EMAIL")
	edgarClient := edgar.New(agentName, agentEmail, time.Second*10)

	formType := "13F-HR"
	resCount := 100
	params := &emodels.LatestFilingsParams{
		Action: emodels.GetCurrent,
		Type:   &formType,
		Count:  &resCount,
		Output: emodels.Atom,
	}
	res, err := edgarClient.FetchLatestFiling(ctx, params)
	if err != nil {
		fmt.Fprintf(stderr, "error getting latest filings: %v\n", err)
	}

	for _, entry := range res.Entries {
		fmt.Fprintf(stdout, "%#v\n\n", entry)
		fmt.Println()
		path, _ := edgarClient.InfotableURLFromHTML(ctx, entry)
		fmt.Fprintf(stdout, "--%s\n", path)
	}

}

func runPolyGrouped(ctx context.Context, getenv func(string) string, stdout, stderr io.Writer) {
	polyClient := polygon.New(getenv("POLYGON_API_KEY"), time.Second*10)
	params := &models.GroupedDailyParams{
		Date: models.Date(time.Date(2024, 12, 12, 0, 0, 0, 0, time.UTC)),
	}
	res, err := polyClient.GroupedDailyBars(ctx, params)
	if err != nil {
		fmt.Fprintf(stderr, "Error happened here\n")
	}
	fmt.Fprintf(stdout, "%+v\n", res)
}

func run(ctx context.Context, getenv func(string) string, stdout, stderr io.Writer) error {
	//dburl := getenv("DB_URL")
	//serveport := getenv("PORT")
	//fmt.Fprintf(stdout, "env variables - dburl: %s, serveport: %s\n", dburl, serveport)

	//runPolyGrouped(ctx, getenv, stdout, stderr)
	//runEdgarFacts(ctx, getenv, stdout, stderr)
	runEdgarFilings(ctx, getenv, stdout, stderr)
	//runEdgarTickers(ctx, getenv, stdout, stderr)
	return nil
}

func main() {
	ctx := context.Background()
	godotenv.Load()
	if err := run(ctx, os.Getenv, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
