package time

import (
	"fmt"
	"testing"
)

func Test_LocalTime(*testing.T) {
	l := NewLocalTime()
	fmt.Printf("%s\n", l.String())
}

func Test_Sub(*testing.T) {
	start := "2016-09-12T08:51:02Z"

	fmt.Println(DurationFromUTC(start))

}
