package events

import (
	"time"
)

func parseTime(s string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t, nil
	}
	return time.Time{}, &time.ParseError{Layout: time.RFC3339, Value: s, LayoutElem: "", ValueElem: ""}
}
