package main

import (
	"context"
	"database/sql"
	"fmt"

	//"fmt"
	//	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jamesonhm/fingator/internal/database"
	//"github.com/jamesonhm/fingator/internal/openfigi"
	//"github.com/jamesonhm/fingator/internal/polygon"
	edgar "github.com/jamesonhm/fingator/internal/sec"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const (
	DaysOHLCVHistory = 3
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Scheduler and Tasks
	s, _ := gocron.NewScheduler(
		gocron.WithLogger(logger),
	)
	defer func() { _ = s.Shutdown() }()

	// defer functions are processed LIFO, context cancel must run before scheduler shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Env Variable Load
	godotenv.Load()
	dburl := mustEnv("DB_URL")
	logger.LogAttrs(ctx, slog.LevelInfo, "Env vars", slog.String("DB_URL", dburl))

	db := must(sql.Open("postgres", dburl))
	defer db.Close()

	err := setup(ctx, db, logger)
	if err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "Setup function failed", slog.Any("Error", err))
		panic(err)
	}

	dbq := database.New(db)

	// API Clients
	//polyClient := polygon.New(
	//	mustEnv("POLYGON_API_KEY"),
	//	time.Second*10,
	//	time.Minute,
	//	5,
	//)

	edgarClient := edgar.New(
		mustEnv("EDGAR_COMPANY_NAME"),
		mustEnv("EDGAR_COMPANY_EMAIL"),
		time.Second*10,
		time.Second,
		10,
	)

	//figiClient := openfigi.New(
	//	os.Getenv("OPENFIGI_API_KEY"),
	//	time.Second*10,
	//)

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
		gocron.WithEventListeners(
			gocron.AfterJobRuns(
				func(jobID uuid.UUID, jobName string) {
					runEdgarFacts(ctx, dbq, edgarClient, logger)
				},
			),
		),
	)

	stmts, err := annualStatements(ctx, dbq, logger, 320193)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	for _, stmt := range stmts {
		fmt.Printf("%+v\n\n", stmt)
	}

	// OHLCV from polygon, weekday-ly
	// TODO: update to CronJob running after close of weekdays
	//_, err = s.NewJob(
	//	gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(5, 0, 0))),
	//	gocron.NewTask(runPolyGrouped, ctx, dbq, polyClient, DaysOHLCVHistory, logger),
	//	gocron.WithContext(ctx),
	//	gocron.WithStartAt(gocron.WithStartImmediately()),
	//)

	//runEdgarCompanyFilings(ctx, dbq, edgarClient, logger)

	//_, err = s.NewJob(
	//	gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(6, 0, 0))),
	//	gocron.NewTask(runEdgarFilings, ctx, dbq, edgarClient, logger),
	//	gocron.WithContext(ctx),
	//	//gocron.WithStartAt(gocron.WithStartImmediately()),
	//	gocron.WithEventListeners(
	//		gocron.AfterJobRuns(
	//			func(jobID uuid.UUID, jobName string) {
	//				runOpenFigiCusips(ctx, dbq, figiClient, logger)
	//			},
	//		),
	//	),
	//)

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

func Xmain() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	godotenv.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	edgarClient := edgar.New(
		mustEnv("EDGAR_COMPANY_NAME"),
		mustEnv("EDGAR_COMPANY_EMAIL"),
		time.Second*10,
		time.Second,
		10,
	)

	runEdgar10k(ctx, edgarClient, logger)

}
