package dating_test

import (
	"math"
	"testing"
	"time"

	"github.com/chaimleib/zmanspec/dating"
)

func TestJulian(t *testing.T) {
	cases := []struct {
		Date string
		Want float64
	}{
		{"1995-10-09T12:00:00Z", 2_450_000},
		{"2023-02-24T12:00:00Z", 2_460_000},
		{"2000-01-01T12:00:00Z", 2_451_545},
		{"2000-01-02T00:00:00Z", 2_451_545.5},
		{"2013-01-01T00:30:00Z", 2_456_293.520_833_33},
		{"2026-01-19T18:02:03Z", 2_461_060.251_423_61},
	}
	for _, c := range cases {
		t.Run(c.Date, func(t *testing.T) {
			d, err := time.Parse(time.RFC3339, c.Date)
			if err != nil {
				t.Fatal(err)
			}
			got := dating.Julian(d)
			if math.Abs(c.Want-got) > 0.000_000_1 {
				t.Errorf("from %s want %f, got %f", c.Date, c.Want, got)
			}
		})
	}
}

func TestFromJulian(t *testing.T) {
	cases := []struct {
		J    float64
		Want string
	}{
		{2_450_000, "1995-10-09T12:00:00Z"},
		{2_460_000, "2023-02-24T12:00:00Z"},
		{2_451_545, "2000-01-01T12:00:00Z"},
		{2_451_545.5, "2000-01-02T00:00:00Z"},
		{2_456_293.520_833_33, "2013-01-01T00:30:00Z"},
		{2_461_060.251_423_61, "2026-01-19T18:02:03Z"},
	}
	for _, c := range cases {
		t.Run(c.Want, func(t *testing.T) {
			d, err := time.Parse(time.RFC3339, c.Want)
			if err != nil {
				t.Fatal(err)
			}
			got := dating.FromJulian(c.J)
			diff := d.Sub(got)
			absDiff := diff
			if diff < 0 {
				absDiff = -absDiff
			}
			if absDiff > time.Millisecond {
				t.Errorf(
					"from %f want %s, got %s which is %s different",
					c.J,
					c.Want,
					got,
					diff,
				)
			}
		})
	}
}
