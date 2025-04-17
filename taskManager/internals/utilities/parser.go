package utilities

import (
	"database/sql/driver"
	"strings"
	"time"
)

type JSONTime time.Time

// method signature
func (t *JSONTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" {
		return nil
	}

	parsed, err := time.Parse(layout, string(s))
	if err != nil {
		return err
	}

	*t = JSONTime(parsed)
	return nil
}

// Value makes JSONTime implement the driver.Valuer interface. it's a method signature
func (t JSONTime) Value() (driver.Value, error) {
	return time.Time(t), nil
}
