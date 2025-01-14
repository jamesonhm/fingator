package main

import (
	"math"
	"time"
)

type DateIter struct {
	Date    time.Time
	days    int
	minDate *time.Time
	maxDate *time.Time
	today   time.Time
}

func NewDateIter(days int, minDate, maxDate *time.Time, today time.Time) *DateIter {
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
	if i.maxDate == nil {
		next = previousWeekday(i.today)
		i.maxDate = &next
		i.Date = midnight(next)
		return true
	}
	if i.today.Sub(*i.maxDate) < time.Hour*24 {
		return false
	}
	next = nextWeekday(*i.maxDate)

	if next.After(i.today) {
		return false
	}
	i.maxDate = &next
	i.Date = midnight(next)
	return true
}

func (i *DateIter) nextMin() bool {
	var next time.Time
	if i.minDate == nil {
		next = previousWeekday(*i.maxDate)
		i.minDate = &next
		i.Date = midnight(next)
		return true
	}
	if weekdaysBetween(*i.minDate, *i.maxDate) > i.days {
		return false
	}
	next = previousWeekday(*i.minDate)
	i.minDate = &next
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

func previousWeekday(d time.Time) time.Time {
	if d.Weekday() == 1 {
		return d.AddDate(0, 0, -3)
	}
	return d.AddDate(0, 0, -1)
}

func nextWeekday(d time.Time) time.Time {
	if d.Weekday() == 5 {
		return d.AddDate(0, 0, 3)
	}
	return d.AddDate(0, 0, 1)
}
