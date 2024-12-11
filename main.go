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
	"github.com/joho/godotenv"
)

func run(ctx context.Context, getenv func(string) string, stdout, stderr io.Writer) error {
	//dburl := getenv("DB_URL")
	//serveport := getenv("PORT")
	//fmt.Fprintf(stdout, "env variables - dburl: %s, serveport: %s\n", dburl, serveport)

	polyClient := polygon.New(getenv("POLYGON_API_KEY"), time.Second*10)

	tType := "CS"
	params := &models.ListTickersParams{
		Type: &tType,
	}
	iter := polyClient.ListTickers(ctx, params)
	for iter.Next() {
		fmt.Fprintf(stdout, "%+v\n", iter.Item())
	}
	if iter.Err() != nil {
		fmt.Fprintf(stdout, "%v\n", iter.Err())
	}
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
