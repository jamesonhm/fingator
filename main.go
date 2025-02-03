package main

import (
	//"database/sql"
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jamesonhm/fingator/internal/database"
	"github.com/jamesonhm/fingator/internal/openfigi"
	"github.com/jamesonhm/fingator/internal/polygon"
	//"github.com/jamesonhm/fingator/internal/rate"

	//edgar "github.com/jamesonhm/fingator/internal/sec"
	"github.com/go-co-op/gocron/v2"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// TODO: Explore the "Must" pattern for env variables and others

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

//func main() {
//	ctx := context.Background()
//	godotenv.Load()
//	if err := run(ctx, os.Getenv, os.Stdout, os.Stderr); err != nil {
//		fmt.Fprintf(os.Stderr, "%s\n", err)
//		os.Exit(1)
//	}
//}

// func run(ctx context.Context, getenv func(string) string, stdout, stderr io.Writer) error {
func main() {
	ctx, cancel := context.WithCancel(context.Background())

	// Env Variable Load
	godotenv.Load()
	dburl := mustEnv("DB_URL")
	db := must(sql.Open("postgres", dburl))
	defer db.Close()
	dbq := database.New(db)

	// API Clients
	polyClient := polygon.New(mustEnv("POLYGON_API_KEY"), time.Second*10, time.Minute, 5)

	// Scheduler and Tasks
	s, _ := gocron.NewScheduler()
	defer func() { _ = s.Shutdown() }()

	//rl := rate.New(time.Second, 2)
	//_, _ = s.NewJob(
	//	gocron.DurationJob(10*time.Second),
	//	gocron.NewTask(
	//		func(ctx context.Context, a string, b int) {
	//			for i := range b {
	//				<-rl.Throttle
	//				fmt.Println(time.Now(), "a:", a, "b:", i)
	//			}
	//		},
	//		"hello",
	//		8,
	//	),
	//	gocron.WithContext(ctx),
	//)

	polyGrouped, err := s.NewJob(
		gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(5, 0, 0))),
		gocron.NewTask(runPolyGrouped, ctx, dbq, polyClient, 7, os.Stdout, os.Stderr),
		gocron.WithContext(ctx),
		gocron.WithStartAt(gocron.WithStartImmediately()),
	)
	if err != nil {
		fmt.Printf("error creating job: %v\n", err)
		return
	}

	fmt.Println(polyGrouped.ID())
	s.Start()

	// Handle interrupt signals (ctrl-c) for graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Block until an interrupt signal is recieved
	<-c
	fmt.Println("\nMain func exit...")
	cancel()

}

func mustEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("Env variable %s required", key))
	}
	return val
}

func must[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}
	return val
}
