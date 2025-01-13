package main

import (
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
)

func TestDateIterMax(t *testing.T) {
	minDate := time.Date(2025, 1, 8, 14, 0, 0, 0, time.Local)
	maxDate := time.Date(2025, 1, 9, 14, 0, 0, 0, time.Local)
	today := time.Date(2025, 1, 13, 6, 0, 0, 0, time.Local)
	di := NewDateIter(5, minDate, maxDate, today)

	di.Next()

	expected := time.Date(2025, 1, 10, 0, 0, 0, 0, time.Local)
	actual := di.Date
	assert.Equal(t, actual, expected)
}

func TestDateIterMaxWknd(t *testing.T) {
	minDate := time.Date(2025, 1, 2, 14, 0, 0, 0, time.Local)
	maxDate := time.Date(2025, 1, 3, 14, 0, 0, 0, time.Local)
	today := time.Date(2025, 1, 13, 6, 0, 0, 0, time.Local)
	di := NewDateIter(5, minDate, maxDate, today)

	di.Next()

	expected := time.Date(2025, 1, 6, 0, 0, 0, 0, time.Local)
	actual := di.Date
	assert.Equal(t, actual, expected)
}

func TestDateIterMax2(t *testing.T) {
	minDate := time.Date(2025, 1, 2, 14, 0, 0, 0, time.Local)
	maxDate := time.Date(2025, 1, 3, 14, 0, 0, 0, time.Local)
	today := time.Date(2025, 1, 13, 6, 0, 0, 0, time.Local)
	di := NewDateIter(5, minDate, maxDate, today)

	di.Next()
	di.Next()

	expected := time.Date(2025, 1, 7, 0, 0, 0, 0, time.Local)
	actual := di.Date
	assert.Equal(t, actual, expected)
}

func TestDateIterMaxPast(t *testing.T) {
	minDate := time.Date(2025, 1, 9, 14, 0, 0, 0, time.Local)
	maxDate := time.Date(2025, 1, 10, 14, 0, 0, 0, time.Local)
	today := time.Date(2025, 1, 13, 6, 0, 0, 0, time.Local)
	di := NewDateIter(5, minDate, maxDate, today)

	new := di.Next()
	assert.Equal(t, new, false)
}
