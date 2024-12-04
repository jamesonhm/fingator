package main

import (
	//"database/sql"
	"context"
	"fmt"
	"io"
	//"log"
	"os"

	"github.com/joho/godotenv"
)

func run(ctx context.Context, getenv func(string) string, stdout, stderr io.Writer) error {
	dburl := getenv("DB_URL")
	serveport := ("PORT")
	fmt.Fprintf(stdout, "env variables - dburl: %s, serveport: %s\n", dburl, serveport)
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
