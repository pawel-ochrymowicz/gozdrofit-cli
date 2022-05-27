package gozdrofitapi

import (
	"encoding/json"
	"time"
)

const DateFormat = "2006-01-02"
const DateTimeFormat = "2006-01-02T15:04:05"

type Date struct {
	time.Time
}

func (Date *Date) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	t, _ := time.Parse(DateFormat, s)
	Date.Time = t
	return nil
}

func (Date Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(Date.Time.Format(DateFormat))
}

type DateTime struct {
	time.Time
}

func (DateTime *DateTime) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	t, _ := time.Parse(DateTimeFormat, s)
	DateTime.Time = t
	return nil
}

func (DateTime DateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(DateTime.Time.Format(DateTimeFormat))
}
