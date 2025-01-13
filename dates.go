package main

import "time"

type DateIter struct {
	Date    time.Time
	days    time.Duration
	minDate time.Time
	maxDate time.Time
	today   time.Time
}

func NewDateIter(days int64, minDate, maxDate, today time.Time) *DateIter {
	dur := time.Duration(int64(time.Hour) * days * 24)
	return &DateIter{
		days:    dur,
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
	if i.maxDate.Sub(i.minDate) > i.days {
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
	offset := -int(start.Weekday())
	start = start.AddDate(0, 0, -int(start.Weekday()))

	offset += int(end.Weekday())
}
