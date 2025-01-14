package main

import (
	"math"
	"time"
)

type DateIter struct {
	Date    time.Time
	days    int
	minDate time.Time
	maxDate time.Time
	today   time.Time
}

func NewDateIter(days int, minDate, maxDate, today time.Time) *DateIter {
	return &DateIter{
		days:    days,
		minDate: minDate,
		maxDate: maxDate,
		today:   midnight(today),
	}
}

func (i *DateIter) Next() bool {
	if i.nextMax() {
		return true
	} else if i.nextMin() {
		return true
	}
	return false
}

func (i *DateIter) nextMax() bool {
	var next time.Time
	if i.today.Sub(i.maxDate) < time.Hour*24 {
		return false
	}
	if i.maxDate.Weekday() == 5 {
		next = i.maxDate.AddDate(0, 0, 3)
	} else {
		next = i.maxDate.AddDate(0, 0, 1)
	}

	if next.After(i.today) {
		return false
	}
	i.maxDate = next
	i.Date = midnight(next)
	return true
}

func (i *DateIter) nextMin() bool {
	var next time.Time
	if weekdaysBetween(i.minDate, i.maxDate) > i.days {
		return false
	}
	if i.minDate.Weekday() == 1 {
		next = i.minDate.AddDate(0, 0, -3)
	} else {
		next = i.minDate.AddDate(0, 0, -1)
	}
	i.minDate = next
	i.Date = midnight(next)
	return true
}

func midnight(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

func weekdaysBetween(start, end time.Time) int {
	offset := int(end.Weekday()) - int(start.Weekday())
	if end.Weekday() == time.Sunday {
		offset++
	}
	start = start.AddDate(0, 0, -int(start.Weekday()))
	end = end.AddDate(0, 0, -int(end.Weekday()))
	diff := end.Sub(start).Truncate(time.Hour * 24)
	weeks := float64((diff.Hours() / 24) / 7)
	return int(math.Round(weeks)*5) + offset
}
