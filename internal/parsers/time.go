package parsers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type TimeParser struct {
	layouts map[string]string
}

func NewTimeParser(layouts map[string]string) *TimeParser {
	return &TimeParser{layouts: layouts}
}

func (p *TimeParser) ParseTime(v string) (*TimeResult, error) {
	tryiedAs := make([]string, 0, len(p.layouts))
	// Try parsing with each layout
	for name, layout := range p.layouts {
		tryiedAs = append(tryiedAs, name)
		if t, err := time.Parse(layout, v); err == nil {
			return &TimeResult{
				TriedAs:  tryiedAs,
				ParsedAs: name,
				Time:     t,
			}, nil
		}
	}
	// Try parsing as Unix timestamp
	if trimmed := strings.TrimSpace(v); trimmed != "" && trimmed != "0" {
		tryiedAs = append(tryiedAs, UnixMilli)
		if n, err := strconv.ParseInt(trimmed, 10, 64); err == nil {
			return &TimeResult{
				TriedAs:  tryiedAs,
				ParsedAs: UnixMilli,
				Time:     time.UnixMilli(n),
			}, nil
		}
	}
	return &TimeResult{
		TriedAs: tryiedAs, // TODO: Maybe error in result? This is not great.
	}, fmt.Errorf("%w: '%s'", errInvalidTimeFormat, v)
}

type TimeResult struct {
	TriedAs  []string
	ParsedAs string
	Time     time.Time
}

var errInvalidTimeFormat = errors.New("invalid time format")

var DefaultLayouts = map[string]string{
	Layout:      time.Layout,
	ANSIC:       time.ANSIC,
	UnixDate:    time.UnixDate,
	RubyDate:    time.RubyDate,
	RFC822:      time.RFC822,
	RFC822Z:     time.RFC822Z,
	RFC850:      time.RFC850,
	RFC1123:     time.RFC1123,
	RFC1123Z:    time.RFC1123Z,
	RFC3339:     time.RFC3339,
	RFC3339Nano: time.RFC3339Nano,
	Kitchen:     time.Kitchen,
	Stamp:       time.Stamp,
	StampMilli:  time.StampMilli,
	StampMicro:  time.StampMicro,
	StampNano:   time.StampNano,
	DateTime:    time.DateTime,
	DateOnly:    time.DateOnly,
	TimeOnly:    time.TimeOnly,
}

const (
	Layout      = "Layout"
	ANSIC       = "ANSIC"
	UnixDate    = "UnixDate"
	RubyDate    = "RubyDate"
	RFC822      = "RFC822"
	RFC822Z     = "RFC822Z"
	RFC850      = "RFC850"
	RFC1123     = "RFC1123"
	RFC1123Z    = "RFC1123Z"
	RFC3339     = "RFC3339"
	RFC3339Nano = "RFC3339Nano"
	Kitchen     = "Kitchen"
	Stamp       = "Stamp"
	StampMilli  = "StampMilli"
	StampMicro  = "StampMicro"
	StampNano   = "StampNano"
	DateTime    = "DateTime"
	DateOnly    = "DateOnly"
	TimeOnly    = "TimeOnly"
	UnixMilli   = "UnixMilli"
)
