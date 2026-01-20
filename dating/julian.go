package dating

import (
	"time"
)

// Julian takes a Gregorian date and returns its Julian day number.
// https://en.wikipedia.org/wiki/Julian_day
func Julian(date time.Time) float64 {
	date = date.UTC()
	y, m, d := date.Date()
	mm := (int(m) - 14) / 12
	j := 1461*(y+4800+mm)/4 +
		367*(int(m)-2-12*mm)/12 -
		(3*(y+4900+mm)/100)/4 +
		d - 32075
	return float64(j) - 0.5 +
		(float64(date.Hour()) / 24.0) +
		(float64(date.Minute()) / 24.0 / 60.0) +
		(float64(date.Second()) / 24.0 / 3600.0) +
		(float64(date.Nanosecond()) / 24.0 / 3600.0 / 1_000_000_000.0)
}

// FromJulian converts a Julian day number into a Gregorian date.
// https://en.wikipedia.org/wiki/Julian_day
func FromJulian(j float64) time.Time {
	jj := int(j + 0.5)
	f := jj + 1401 + (((4*jj+274277)/146097)*3)/4 - 38
	e := 4*f + 3
	g := e % 1461 / 4
	h := 5*g + 2
	d := (h%153)/5 + 1
	m := (h/153+2)%12 + 1
	y := (e / 1461) - 4716 + (12+2-m)/12
	frac := j + 0.5 - float64(int(j+0.5))
	hh := 24.0 * frac
	mm := (hh - float64(int(hh))) * 60.0
	ss := (mm - float64(int(mm))) * 60.0
	ns := (ss - float64(int(ss))) * 1_000_000_000
	return time.Date(
		y,
		time.Month(m),
		d,
		int(hh),
		int(mm),
		int(ss),
		int(ns),
		time.UTC,
	)
}
