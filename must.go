package main

import (
	"fmt"
	"os"
)

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
