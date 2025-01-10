package main

import "time"

type DateIter struct {
	Date    time.Time
	days    time.Duration
	minDate time.Time
	maxDate time.Time
	today   time.Time
}

func NewDateIter(days time.Duration, minDate, maxDate time.Time) *DateIter {
	return &DateIter{
		days:    days,
		minDate: minDate,
		maxDate: maxDate,
		today:   midnight(time.Now()),
	}
}

func (i *DateIter) Next() bool {
	if i.nextMax() {
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
		next = i.maxDate.Add(time.Hour * 72)
	} else {
		next = i.maxDate.Add(time.Hour * 24)
	}

	if next.After(i.today) {
		return false
	}
	i.Date = midnight(next)
	return true
}

func midnight(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}
