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
	}
	for _, c := range cases {
		t.Run(c.Want, func(t *testing.T) {
			d, err := time.Parse(time.RFC3339, c.Want)
			if err != nil {
				t.Fatal(err)
			}
			got := dating.FromJulian(c.J)
			if !d.Equal(got) {
				t.Errorf("from %f want %s, got %s", c.J, c.Want, got)
			}
		})
	}
}
