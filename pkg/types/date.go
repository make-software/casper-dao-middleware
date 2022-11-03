package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

var ErrInvalidDateFormat = errors.New("invalid date format provided")

// Date used to store datestamps
type Date time.Time

func NewDateFromTime(date time.Time) Date {
	return Date(date.Truncate(24 * time.Hour))
}

func ParseDateFromString(date string) (Date, error) {
	parsed, err := time.Parse("2006-01-02", date)
	if err != nil {
		return Date{}, ErrInvalidDateFormat
	}
	return Date(parsed), nil
}

func ParseDateFromTime(date time.Time) (Date, error) {
	if err := IsDate(date); err != nil {
		return Date{}, err
	}
	return Date(date), nil
}

func (d *Date) AddDate(years int, months int, days int) Date {
	dateTime := time.Time(*d).AddDate(years, months, days)
	return Date(dateTime)
}

// MarshalJSON convert Date in format 2006-01-02
func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(d).Format("2006-01-02"))
}

func (d Date) Time() time.Time {
	return time.Time(d)
}

// Value rewrite behaviour for inserting Date to db
func (d Date) Value() (driver.Value, error) {
	return time.Time(d), nil
}

func (d *Date) Scan(value interface{}) error {
	if value == nil {
		return errors.New("nil date")
	}

	bv, err := driver.DefaultParameterConverter.ConvertValue(value)
	if err != nil {
		return errors.New("invalid date value")
	}

	raw, ok := bv.(time.Time)
	if !ok {
		return errors.New("invalid date value")
	}

	date, err := ParseDateFromTime(raw)
	if err != nil {
		return err
	}
	*d = date
	return nil
}

func IsDate(date time.Time) error {
	if date.IsZero() {
		return errors.New("zero date provided")
	}

	hour, min, sec := date.Clock()
	if hour != 0 || min != 0 || sec != 0 {
		return ErrInvalidDateFormat
	}
	return nil
}

func IsListOfDates(dates []time.Time) error {
	for _, date := range dates {
		if err := IsDate(date); err != nil {
			return err
		}
	}
	return nil
}
