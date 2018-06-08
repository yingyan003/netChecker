package day

import (
	"time"
)

type Day struct {
	Day time.Time
}

type Days []Day

func NewToday() *Day {
	d := new(Day)
	d.Day = time.Now()
	return d
}

func NewDay(str string) *Day {
	d := new(Day)
	d.Day, _ = time.Parse("2006-01-02", str)
	return d
}

func (d *Day) String() string {
	return d.Day.Format("2006-01-02")
}

func (d *Day) GetYesterday() *Day {
	yesterday := new(Day)
	yesterday.Day = d.Day.AddDate(0, 0, -1)
	return yesterday
}

func (d *Day) GetTomorrow() *Day {
	tomorrow := new(Day)
	tomorrow.Day = d.Day.AddDate(0, 0, 1)
	return tomorrow
}

func (d *Day) GetLastDays(n int) Days {
	lastDays := make([]Day, 0)
	for i := 0; i < n; i++ {
		last := new(Day)
		last.Day = d.Day.AddDate(0, 0, -i)
		lastDays = append(lastDays, *last)
	}

	return lastDays
}

func (d *Day) GetLastDaysString(n int) []string {
	lastDays := d.GetLastDays(n)
	lastDaysString := make([]string, 0)
	for _, last := range lastDays {
		lastDaysString = append(lastDaysString, last.String())
	}
	return lastDaysString
}

func (ds Days) Len() int {
	return len(ds)
}

func (ds Days) Swap(i, j int) {
	ds[i], ds[j] = ds[j], ds[i]
}

func (ds Days) Less(i, j int) bool {
	return ds[i].String() < ds[j].String()
}
