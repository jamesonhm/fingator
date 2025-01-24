package models

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type NumericCIK int

func (n NumericCIK) Pad() string {
	return fmt.Sprintf("%010d", n)
}

// Json Decoding Function
type DecFunc func(r *http.Response, v any) error

// Date is a short date without a time component of the format: "2006-01-02"
type Date time.Time

// PathFormat used to string format for use as a path parameter
func (d Date) PathFormat() string {
	return time.Time(d).Format("2006-01-02")
}

func (d *Date) UnmarshalJSON(data []byte) error {
	unquoteData, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}
	t, err := time.Parse("2006-01-02", unquoteData)
	if err != nil {
		return err
	}
	*d = Date(t)
	return nil
}

// Stringer Interface for print format
func (d Date) String() string {
	return fmt.Sprintf("%v", time.Time(d).Format(time.DateOnly))
}

// Time is a date-time
type Time time.Time

//func (t *Time) UnmarshalXML(data []byte) error {
//	unquoteData, err := strconv.Unquote(string(data))
//	if err != nil {
//		return err
//	}
//
//	if parsedTime, err := time.Parse("2006-01-02T15:04:05.000-0700", unquoteData); err == nil {
//		*t = Time(parsedTime)
//		return nil
//	}
//
//	if parsedTime, err := time.Parse("2006-01-02T15:04:05-07:00", unquoteData); err == nil {
//		*t = Time(parsedTime)
//		return nil
//	}
//
//	if parsedTime, err := time.Parse("2006-01-02T15:04:05.000Z", unquoteData); err == nil {
//		*t = Time(parsedTime)
//		return nil
//	}
//
//	if parsedTime, err := time.Parse("2006-01-02T15:04:05Z", unquoteData); err != nil {
//		return err
//	} else {
//		*t = Time(parsedTime)
//	}
//
//	return nil
//}

type Action string

const (
	GetCurrent Action = "getcurrent"
	GetCompany Action = "getcompany"
)

type Ownership string

const (
	Include Ownership = "include"
	Exclude Ownership = "exclude"
	Only    Ownership = "only"
)

type Output string

const (
	Atom Output = "atom"
)
