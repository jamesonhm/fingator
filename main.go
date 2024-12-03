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
	return nil
}

func main() {
	ctx := context.Background()
	if err := run(ctx, getenv, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func getenv(key string) string {

	godotenv.Load()
	switch key {
	case "DB_URL":
		return os.Getenv("DB_URL")
	case "PORT":
		return os.Getenv("PORT")
	default:
		return ""
		//db, err := sql.Open("postgres", dbURL)
		//if err != nil {
		//	log.Fatal(err)
		//}
	}
}
