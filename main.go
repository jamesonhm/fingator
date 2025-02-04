package main

import (
	//"database/sql"
	"context"
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jamesonhm/fingator/internal/database"
	"github.com/jamesonhm/fingator/internal/openfigi"
	"github.com/jamesonhm/fingator/internal/polygon"

	//"github.com/jamesonhm/fingator/internal/rate"

	"github.com/go-co-op/gocron/v2"
	edgar "github.com/jamesonhm/fingator/internal/sec"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func Xrun(ctx context.Context, getenv func(string) string, stdout, stderr io.Writer) error {
	dburl := getenv("DB_URL")
	db, err := sql.Open("postgres", dburl)
	if err != nil {
		return fmt.Errorf("unable to connect to db: %v", err)
	}
	defer db.Close()
	//dbq := database.New(db)

	//serveport := getenv("PORT")
	//fmt.Fprintf(stdout, "env variables - dburl: %s, serveport: %s\n", dburl, serveport)

	//polyClient := polygon.New(getenv("POLYGON_API_KEY"), time.Second*10, 0.083)
	//runPolyGrouped(ctx, dbq, polyClient, stdout, stderr)

	//agentName := getenv("EDGAR_COMPANY_NAME")
	//agentEmail := getenv("EDGAR_COMPANY_EMAIL")
	//edgarClient := edgar.New(agentName, agentEmail, time.Second*10, 10)

	//runEdgarTickers(ctx, dbq, edgarClient, stdout, stderr)

	//runEdgarFacts(ctx, dbq, edgarClient, stdout, stderr)

	//runEdgarFilings(ctx, getenv, stdout, stderr)
	//fmt.Fprintf(stdout, "==============================================\n")
	//runEdgarCompanyFilings(ctx, getenv, stdout, stderr)

	figiClient := openfigi.New(getenv("OPENFIGI_API_KEY"), time.Second*10, 4)
	runOpenFigiCusips(ctx, figiClient, stdout, stderr)

	return nil
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Env Variable Load
	godotenv.Load()
	dburl := mustEnv("DB_URL")
	db := must(sql.Open("postgres", dburl))
	defer db.Close()
	dbq := database.New(db)

	// API Clients
	polyClient := polygon.New(
		mustEnv("POLYGON_API_KEY"),
		time.Second*10,
		time.Minute,
		5,
	)

	edgarClient := edgar.New(
		mustEnv("EDGAR_COMPANY_NAME"),
		mustEnv("EDGAR_COMPANY_EMAIL"),
		time.Second*10,
		time.Second,
		10,
	)

	// Scheduler and Tasks
	s, _ := gocron.NewScheduler(
		gocron.WithLogger(logger),
	)
	defer func() { _ = s.Shutdown() }()

	// defer functions are processed LIFO, context cancel must run before scheduler shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// OHLCV from polygon, weekday-ly
	// TODO: update to CronJob running after close of weekdays
	_, err := s.NewJob(
		gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(5, 0, 0))),
		gocron.NewTask(runPolyGrouped, ctx, dbq, polyClient, 7, logger),
		gocron.WithContext(ctx),
		gocron.WithStartAt(gocron.WithStartImmediately()),
	)

	// CIK/Ticker/Exchange from Edgar, monthly
	_, err = s.NewJob(
		gocron.MonthlyJob(
			1,
			gocron.NewDaysOfTheMonth(1),
			gocron.NewAtTimes(gocron.NewAtTime(12, 0, 0)),
		),
		gocron.NewTask(runEdgarTickers, ctx, dbq, edgarClient, logger),
		gocron.WithContext(ctx),
		gocron.WithStartAt(gocron.WithStartImmediately()),
	)
	if err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "error creating job", slog.Any("Error:", err))
		return
	}

	s.Start()

	// Handle interrupt signals (ctrl-c) for graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Block until an interrupt signal is recieved
	<-c
	logger.InfoContext(ctx, "Main func exit...")

}
