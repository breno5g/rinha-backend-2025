package utils

import "time"

func ParseTime(param string) *time.Time {
	if param == "" {
		return nil
	}
	if t, err := time.Parse(time.RFC3339, param); err == nil {
		return &t
	}
	return nil
}
