package zmanspec

import (
	"fmt"
	"strings"
)

type ColumnParseError struct {
	Err error
	Col int
	S   string
}

func (pe ColumnParseError) Error() string {
	if pe.Col < 1 {
		return fmt.Sprintf("ColumnParseError: invalid Col %d; %v", pe.Col, pe.Err)
	}

	var s, marker strings.Builder
	indent := strings.Repeat(" ", 8)
	var i int
	for ; i < pe.Col-1; i++ {
		switch pe.S[i] {
		case '\t':
			s.WriteString(indent)
			marker.WriteString(indent)
		default:
			cs := fmt.Sprintf("%c", pe.S[i])
			s.WriteString(cs)
			marker.WriteString(strings.Repeat(" ", len(cs)))
		}
	}
	marker.WriteRune('^')
	for ; i < len(pe.S); i++ {
		fmt.Fprintf(&s, "%c", pe.S[i])
	}

	return fmt.Sprintf("%v\n  %s\n  %s", pe.Err, &s, &marker)
}

func (pe ColumnParseError) Unwrap() error {
	return pe.Err
}
