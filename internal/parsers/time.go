package parsers

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

type TimeParser struct {
	layouts []LayoutName
}

func NewTimeParser(layouts ...LayoutName) *TimeParser {
	return &TimeParser{layouts: layouts}
}

func NewTimeParserAll() *TimeParser {
	return NewTimeParser(
		Layout,
		ANSIC,
		UnixDate,
		RubyDate,
		RFC822,
		RFC822Z,
		RFC850,
		RFC1123,
		RFC1123Z,
		RFC3339,
		RFC3339Nano,
		Kitchen,
		Stamp,
		StampMilli,
		StampMicro,
		StampNano,
		DateTime,
		DateOnly,
		TimeOnly,
		UnixMilli,
	)
}

func (p *TimeParser) ParseTime(v string) (*TimeResult, error) {
	tryiedAs := make([]LayoutName, 0, len(p.layouts))
	// Try parsing with each layout
	for _, layout := range p.layouts {
		tryiedAs = append(tryiedAs, layout)
		if t, err := parseLayout(layout, v); err == nil {
			return &TimeResult{
				TriedAs:  tryiedAs,
				ParsedAs: layout,
				Time:     t,
			}, nil
		}
	}
	return nil, fmt.Errorf("%w: '%s'", errInvalidTimeFormat, v)
}

type TimeResult struct {
	TriedAs  []LayoutName
	ParsedAs LayoutName
	Time     time.Time
}

var (
	errUnknownLayout     = errors.New("unknown layout")
	errInvalidTimeFormat = errors.New("invalid time format")
)

type LayoutName string

const (
	Layout      LayoutName = "Layout"
	ANSIC       LayoutName = "ANSIC"
	UnixDate    LayoutName = "UnixDate"
	RubyDate    LayoutName = "RubyDate"
	RFC822      LayoutName = "RFC822"
	RFC822Z     LayoutName = "RFC822Z"
	RFC850      LayoutName = "RFC850"
	RFC1123     LayoutName = "RFC1123"
	RFC1123Z    LayoutName = "RFC1123Z"
	RFC3339     LayoutName = "RFC3339"
	RFC3339Nano LayoutName = "RFC3339Nano"
	Kitchen     LayoutName = "Kitchen"
	Stamp       LayoutName = "Stamp"
	StampMilli  LayoutName = "StampMilli"
	StampMicro  LayoutName = "StampMicro"
	StampNano   LayoutName = "StampNano"
	DateTime    LayoutName = "DateTime"
	DateOnly    LayoutName = "DateOnly"
	TimeOnly    LayoutName = "TimeOnly"
	UnixMilli   LayoutName = "UnixMilli"
)

func parseLayout(layout LayoutName, v string) (time.Time, error) {
	switch layout {
	case Layout:
		return time.Parse(time.Layout, v)
	case ANSIC:
		return time.Parse(time.ANSIC, v)
	case UnixDate:
		return time.Parse(time.UnixDate, v)
	case RubyDate:
		return time.Parse(time.RubyDate, v)
	case RFC822:
		return time.Parse(time.RFC822, v)
	case RFC822Z:
		return time.Parse(time.RFC822Z, v)
	case RFC850:
		return time.Parse(time.RFC850, v)
	case RFC1123:
		return time.Parse(time.RFC1123, v)
	case RFC1123Z:
		return time.Parse(time.RFC1123Z, v)
	case RFC3339:
		return time.Parse(time.RFC3339, v)
	case RFC3339Nano:
		return time.Parse(time.RFC3339Nano, v)
	case Kitchen:
		return time.Parse(time.Kitchen, v)
	case Stamp:
		return time.Parse(time.Stamp, v)
	case StampMilli:
		return time.Parse(time.StampMilli, v)
	case StampMicro:
		return time.Parse(time.StampMicro, v)
	case StampNano:
		return time.Parse(time.StampNano, v)
	case DateTime:
		return time.Parse(time.DateTime, v)
	case DateOnly:
		return time.Parse(time.DateOnly, v)
	case TimeOnly:
		return time.Parse(time.TimeOnly, v)
	case UnixMilli:
		vInt, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return time.Time{}, err
		}
		return time.UnixMilli(vInt), nil
	default:
		return time.Time{}, errUnknownLayout
	}
}
