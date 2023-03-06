package main

import (
	"fmt"
	"strings"
	"time"
)

type Date time.Time

func (d Date) String() string {
	return time.Time(d).Format("2006-01-02")
}


type Dates []Date

func (d Dates) String() string {
	strDates := make([]string, len(d))
	for i, date := range d {
		strDates[i] = date.String()
	}
	return strings.Join(strDates, ",")
}

func parseDate(date string) (time.Time, error) {
	return time.Parse("2006-1-2", date)
}

// NewDates returns a new Dates object. If `dates` is empty, it will return an
// empty Dates object, but no error.
func NewDates(dates string) (Dates, error) {
	dts := Dates{}
	if dates == "" {
		return dts, nil
	}

	dtsOrRanges := strings.Split(dates, ",")
	for _, dateOrRange := range dtsOrRanges {
		rng := strings.Split(dateOrRange, ":")
		if len(rng) == 1 {
			// Single date.
			date, err := parseDate(dateOrRange)
			if err != nil {
				return dts, err
			}
			dts = append(dts, Date(date))

		} else if len(rng) == 2 {
			// Range of dates.
			d1, e1 := parseDate(rng[0])
			if e1 != nil {
				return dts, e1
			}
			d2, e2 := parseDate(rng[1])
			if e2 != nil {
				return dts, e2
			}
			// Fix wrong order.
			for d1.After(d2) {
				d1, d2 = d2, d1
			}
			for d1.Before(d2) {
				dts = append(dts, Date(d1))
				d1 = d1.Add(time.Hour * 24)
			}
			dts = append(dts, Date(d2))

		} else {
			return dts, fmt.Errorf("incorrect range: %s", dates)
		}
	}

	return dts, nil
}
