package zmanspec

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

func leadingInt(s string) (x uint64, rem string, err error) {
	i := 0
	for ; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			break
		}
		if i > 0 && s[0] == '0' {
			return 0, s[i:], errors.New("no digits expected after initial zero")
		}
		if i > 4 {
			// also return rem for the error message's column marker
			return 0, s[i:], errors.New("too many consecutive digits")
		}
		x = x*10 + uint64(c-'0')
	}
	return x, s[i:], nil
}

func ParseDuration(s string) (time.Duration, error) {
	if s == "0" {
		return 0, nil
	}

	orig := s
	columnErr := func(msg string) error {
		return ColumnParseError{
			Err: errors.New("invalid duration: " + msg),
			Col: 1 + len(orig) - len(s),
			S:   orig,
		}
	}

	var d time.Duration
	var neg bool
	if s != "" {
		c := s[0]
		if c == '-' {
			neg = true
			s = s[1:]
		}
	}

	if s == "" {
		return 0, columnErr("unexpected end of string, expected a number")
	}

	units := "hms"
	unitMuls := []time.Duration{
		time.Hour,
		time.Minute,
		time.Second,
	}

	for s != "" {
		var (
			v   uint64
			err error
		)
		// Consume int
		v, s, err = leadingInt(s)
		if err != nil {
			return 0, columnErr(err.Error())
		}

		// Consume unit
		if s == "" {
			return 0, columnErr(fmt.Sprintf(
				"unexpected end of string, expected one of the units %q",
				units,
			))
		}

		unitIdx := strings.IndexByte(units, s[0])
		if unitIdx < 0 {
			return 0, columnErr(fmt.Sprintf("expected one of the units %q", units))
		}
		d += unitMuls[unitIdx] * time.Duration(v)
		units = units[unitIdx+1:]
		unitMuls = unitMuls[unitIdx+1:]
		s = s[1:]
		if units == "" && s != "" {
			return 0, columnErr("no allowed units remaining")
		}

	}

	if neg {
		d = -d
	}
	return d, nil
}
