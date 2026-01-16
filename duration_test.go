package zmanspec_test

import (
	"testing"
	"time"

	"github.com/chaimleib/zmanspec"
)

func TestParseDuration(t *testing.T) {
	cases := []struct {
		Name  string
		Input string
		Want  time.Duration
		Err   string
	}{
		{
			Name: "empty",
			Err: `invalid duration: unexpected end of string, expected a number
  
  ^`,
		},
		{Name: "zero", Input: "0"},
		{
			Name:  "missing unit",
			Input: "1",
			Err: `invalid duration: unexpected end of string, expected one of the units "hms"
  1
   ^`,
		},
		{
			Name:  "invalid unit",
			Input: "1x",
			Err: `invalid duration: expected one of the units "hms"
  1x
   ^`,
		},
		{
			Name:  "invalid number",
			Input: "h",
			Err: `invalid duration: expected non-zero digit
  h
  ^`,
		},
		{
			Name:  "invalid zero hours should be unitless zero",
			Input: "0h",
			Err: `invalid duration: expected non-zero digit
  0h
  ^`,
		},
		{
			Name:  "too many digits",
			Input: "123456s",
			Err: `invalid duration: too many consecutive digits
  123456s
       ^`,
		},
		{
			Name:  "repeated units",
			Input: "6m3m",
			Err: `invalid duration: expected one of the units "s"
  6m3m
     ^`,
		},
		{
			Name:  "repeated zeroed units",
			Input: "0h0h",
			Err: `invalid duration: expected one of the units "ms"
  0h0h
     ^`,
		},
		{
			Name:  "unsorted units",
			Input: "6s3m",
			Err: `invalid duration: no allowed units remaining
  6s3m
    ^`,
		},
		{
			Name:  "unsorted zeroed units",
			Input: "0s0m",
			Err: `invalid duration: no allowed units remaining
  0s0m
    ^`,
		},
		{
			Name:  "negative zero",
			Input: "-0",
			Err: `invalid duration: unexpected end of string, expected one of the units "hms"
  -0
    ^`,
		},
		{
			Name:  "hours",
			Input: "2h",
			Want:  2 * time.Hour,
		},
		{
			Name:  "minutes",
			Input: "3m",
			Want:  3 * time.Minute,
		},
		{
			Name:  "seconds",
			Input: "4s",
			Want:  4 * time.Second,
		},
		{
			Name:  "negative hours",
			Input: "-2h",
			Want:  -2 * time.Hour,
		},
		{
			Name:  "hours and minutes",
			Input: "2h3m",
			Want:  2*time.Hour + 3*time.Minute,
		},
		{
			Name:  "hours and minutes and seconds",
			Input: "2h3m4s",
			Want:  2*time.Hour + 3*time.Minute + 4*time.Second,
		},
		{
			Name:  "hours and seconds",
			Input: "2h4s",
			Want:  2*time.Hour + 4*time.Second,
		},
		{
			Name:  "minutes and seconds",
			Input: "3m4s",
			Want:  3*time.Minute + 4*time.Second,
		},
		{
			Name:  "zero hours",
			Input: "0h",
		},
		{
			Name:  "zero minutes",
			Input: "0m",
		},
		{
			Name:  "zero seconds",
			Input: "0s",
		},
		{
			Name:  "zero hours minutes",
			Input: "0h0m",
		},
		{
			Name:  "zero hms",
			Input: "0h0m0s",
		},
		{
			Name:  "zero hours seconds",
			Input: "0h0s",
		},
		{
			Name:  "zero minutes seconds",
			Input: "0m0s",
		},
		{
			Name:  "too many zeroes",
			Input: "00m",
			Err: `invalid duration: no digits expected after initial zero
  00m
   ^`,
		},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			got, err := zmanspec.ParseDuration(c.Input)
			if err != nil {
				if c.Err != err.Error() {
					t.Errorf("got err:\n%v\nwant err:\n%s", err, c.Err)
				}
			}
			if c.Want != got {
				t.Errorf("want: %s\ngot:  %s", c.Want, got)
			}
		})
	}
}
