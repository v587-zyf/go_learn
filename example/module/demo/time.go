package example

import (
	"fmt"
	"time"
)

func Time() {
	d := getYearWeekByHour(time.Now(), 1, 11, 30, 59)
	fmt.Println(d)
}

func getYearWeekByHour(t time.Time, d, h, m, s int32) int {
	dt := (24 * time.Hour) * time.Duration(d)
	ht := time.Hour * time.Duration(h)
	mt := time.Minute * time.Duration(m)
	st := time.Second * time.Duration(s)

	year, week := t.Add(-(dt + ht + mt + st)).ISOWeek()
	date := year*100 + week
	return date
}
