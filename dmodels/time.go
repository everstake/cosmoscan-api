package dmodels

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Time struct {
	time.Time
}

func NewTime(t time.Time) Time {
	return Time{Time: t}
}

func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(t.Unix(), 10)), nil
}

func (t *Time) UnmarshalJSON(data []byte) error {
	str := string(data)
	str = strings.Trim(str, `"`)
	timestamp, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return err
	}
	t.Time = time.Unix(timestamp, 0)
	return nil
}

const timeFormat = "2006-01-02 15:04:05.999999"

// Scan implements the Scanner interface.
// The value type must be time.Time or string / []byte (formatted time-string),
// otherwise Scan fails.
func (t *Time) Scan(value interface{}) (err error) {
	if value == nil {
		return fmt.Errorf("invalid value")
	}

	switch v := value.(type) {
	case time.Time:
		t.Time = v
		return
	case []byte:
		t.Time, err = parseDateTime(string(v), time.UTC)
		if err != nil {
			return err
		}
	case string:
		t.Time, err = parseDateTime(v, time.UTC)
		if err != nil {
			return err
		}
	}
	return fmt.Errorf("can't convert %T to time.Time", value)
}

// Value implements the driver Valuer interface.
func (t Time) Value() (driver.Value, error) {
	return t.Time, nil
}

func parseDateTime(str string, loc *time.Location) (t time.Time, err error) {
	base := "0000-00-00 00:00:00.0000000"
	switch len(str) {
	case 10, 19, 21, 22, 23, 24, 25, 26: // up to "YYYY-MM-DD HH:MM:SS.MMMMMM"
		if str == base[:len(str)] {
			return
		}
		t, err = time.Parse(timeFormat[:len(str)], str)
	default:
		err = fmt.Errorf("invalid time string: %s", str)
		return
	}

	// Adjust location
	if err == nil && loc != time.UTC {
		y, mo, d := t.Date()
		h, mi, s := t.Clock()
		t, err = time.Date(y, mo, d, h, mi, s, t.Nanosecond(), loc), nil
	}

	return
}

func (t Time) MarshalBinary() ([]byte, error) {
	return t.Time.MarshalBinary()
}

func (t *Time) UnmarshalBinary(data []byte) error {
	return t.Time.UnmarshalBinary(data)
}

// IsZero returns true for null strings (omitempty support)
func (t Time) IsZero() bool {
	return t.Time.IsZero()
}
