package models

import (
	"fmt"
	"strconv"
	"time"
)

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
