package time

import (
	"strconv"
	"time"
)

type LocalTime struct {
}

const DAYINSECONDS = 86400

func NewLocalTime() *LocalTime {
	return &LocalTime{}
}

func (local *LocalTime) String() string {
	t := time.Now()
	return t.Format(time.RFC3339)
}

func DurationFromUTC(t string) string {
	s, _ := time.Parse(time.RFC3339, t)
	start := s.UTC()

	now := time.Now().UTC()

	duration := now.Sub(start)
	hour := duration.Hours()
	min := duration.Minutes()
	sec := duration.Seconds()

	var durationString string

	if hour > 24 {
		durationString = strconv.Itoa(int(hour/24)) + " days"
	} else if hour < 1 {
		if min > 1 {
			durationString = strconv.Itoa(int(min)) + " mins"
		} else {
			durationString = strconv.Itoa(int(sec)) + " secs"
		}
	} else {
		durationString = strconv.Itoa(int(hour)) + " hours"
	}

	return durationString
}
