package main

import (
	//"database/sql"
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"

	"github.com/jamesonhm/fingator/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// TODO: Explore the "Must" pattern for env variables and others

func run(ctx context.Context, getenv func(string) string, stdout, stderr io.Writer) error {
	dburl := getenv("DB_URL")
	db, err := sql.Open("postgres", dburl)
	if err != nil {
		return fmt.Errorf("unable to connect to db: %v", err)
	}
	defer db.Close()
	dbq := database.New(db)
	//serveport := getenv("PORT")
	//fmt.Fprintf(stdout, "env variables - dburl: %s, serveport: %s\n", dburl, serveport)

	runPolyGrouped(ctx, dbq, getenv, stdout, stderr)
	//runEdgarFacts(ctx, getenv, stdout, stderr)
	//runEdgarFilings(ctx, getenv, stdout, stderr)
	//fmt.Fprintf(stdout, "==============================================\n")
	//runEdgarCompanyFilings(ctx, getenv, stdout, stderr)
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
