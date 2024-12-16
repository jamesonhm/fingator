package main

import (
	//"database/sql"
	"context"
	"fmt"
	"io"
	"time"

	//"log"
	"os"

	//"github.com/jamesonhm/fingator/internal/polygon"
	//"github.com/jamesonhm/fingator/internal/polygon/models"
	edgar "github.com/jamesonhm/fingator/internal/sec"
	"github.com/joho/godotenv"
)

func run(ctx context.Context, getenv func(string) string, stdout, stderr io.Writer) error {
	//dburl := getenv("DB_URL")
	//serveport := getenv("PORT")
	//fmt.Fprintf(stdout, "env variables - dburl: %s, serveport: %s\n", dburl, serveport)

	//polyClient := polygon.New(getenv("POLYGON_API_KEY"), time.Second*10)
	agentName := getenv("EDGAR_COMPANY_NAME")
	agentEmail := getenv("EDGAR_COMPANY_EMAIL")
	edgarClient := edgar.New(agentName, agentEmail, time.Second*10)

	_, err := edgarClient.GetCompanyTickers(ctx)
	//params := &models.GroupedDailyParams{
	//	Date: models.Date(time.Date(2024, 12, 12, 0, 0, 0, 0, time.UTC)),
	//}
	//res, err := polyClient.GroupedDailyBars(ctx, params)

	//params := &models.TickerDetailsParams{
	//	Ticker: "AAPL",
	//}
	//res, err := polyClient.GetTickerDetails(ctx, params)
	if err != nil {
		fmt.Fprintf(stderr, "Error happened here\n")
		return err
	}
	//fmt.Fprintf(stdout, "%+v\n", res)
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
