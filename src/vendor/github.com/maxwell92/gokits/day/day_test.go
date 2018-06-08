package day

import (
	"fmt"
	"testing"
)

func TestNewDay(t *testing.T) {
	d := NewDay("2016-10-01")
	fmt.Println("TestNewDay: " + d.String())
}

func TestNewToday(t *testing.T) {
	d := NewToday()
	fmt.Println("TestNewToday: " + d.String())
}

func TestDay_GetYesterday(t *testing.T) {
	d := NewToday()
	yesterday := d.GetYesterday()
	fmt.Println("TestDay_GetYesterday: " + yesterday.String())
}

func TestDay_GetTomorrow(t *testing.T) {
	d := NewToday()
	tomorrow := d.GetTomorrow()
	fmt.Println("TestDay_GetTomorrow: " + tomorrow.String())
}

func TestDay_GetLastDays(t *testing.T) {
	d := NewDay("2016-10-03")
	ds := d.GetLastDays(7)

	fmt.Println("TestDay_GetLastDays")
	for _, day := range ds {
		fmt.Println(day.String())
	}
	fmt.Println()
}

func TestDay_GetLastDaysString(t *testing.T) {
	d := NewToday()
	ds := d.GetLastDaysString(7)

	fmt.Println("TestDay_GetLastDaysString")
	for _, day := range ds {
		fmt.Println(day)
	}
	fmt.Println()
}
