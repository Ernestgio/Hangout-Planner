package types

import (
	"strings"
	"time"

	"github.com/Ernestgio/Hangout-Planner/pkg/shared/constants"
)

type JSONTime time.Time

func (t JSONTime) MarshalJSON() ([]byte, error) {
	if time.Time(t).IsZero() {
		return []byte("null"), nil
	}
	formatted := time.Time(t).Format(constants.DateFormat)
	return []byte(`"` + formatted + `"`), nil
}

func (t *JSONTime) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	if s == "null" || s == "" {
		*t = JSONTime(time.Time{})
		return nil
	}

	parsedTime, err := time.Parse(constants.DateFormat, s)
	if err != nil {
		return err
	}

	*t = JSONTime(parsedTime)
	return nil
}
