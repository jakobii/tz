package chrono

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Format int

const (
	Layout Format = iota
	ANSIC
	UnixDate
	RubyDate
	RFC822
	RFC822Z
	RFC850
	RFC1123
	RFC1123Z
	RFC3339
	RFC3339Nano
	Kitchen
	Stamp
	StampMilli
	StampMicro
	StampNano
	DateTime
	DateOnly
	TimeOnly
	Unix
	UnixMilli
	UnixMicro
)

func (f Format) String() string {
	switch f {
	case Layout:
		return "Layout"
	case ANSIC:
		return "ANSIC"
	case UnixDate:
		return "UnixDate"
	case RubyDate:
		return "RubyDate"
	case RFC822:
		return "RFC822"
	case RFC822Z:
		return "RFC822Z"
	case RFC850:
		return "RFC850"
	case RFC1123:
		return "RFC1123"
	case RFC1123Z:
		return "RFC1123Z"
	case RFC3339:
		return "RFC3339"
	case RFC3339Nano:
		return "RFC3339Nano"
	case Kitchen:
		return "Kitchen"
	case Stamp:
		return "Stamp"
	case StampMilli:
		return "StampMilli"
	case StampMicro:
		return "StampMicro"
	case StampNano:
		return "StampNano"
	case DateTime:
		return "DateTime"
	case DateOnly:
		return "DateOnly"
	case TimeOnly:
		return "TimeOnly"
	case Unix:
		return "Unix"
	case UnixMilli:
		return "UnixMilli"
	case UnixMicro:
		return "UnixMicro"
	default:
		return "Unknown"
	}
}

func ParseFormat(formatStr string) (Format, error) {
	formatStr = strings.ToLower(formatStr)
	switch formatStr {
	case "layout":
		return Layout, nil
	case "ansic":
		return ANSIC, nil
	case "unixdate":
		return UnixDate, nil
	case "rubydate":
		return RubyDate, nil
	case "rfc822":
		return RFC822, nil
	case "rfc822z":
		return RFC822Z, nil
	case "rfc850":
		return RFC850, nil
	case "rfc1123":
		return RFC1123, nil
	case "rfc1123z":
		return RFC1123Z, nil
	case "rfc3339":
		return RFC3339, nil
	case "rfc3339nano":
		return RFC3339Nano, nil
	case "kitchen":
		return Kitchen, nil
	case "stamp":
		return Stamp, nil
	case "stampmilli":
		return StampMilli, nil
	case "stampmicro":
		return StampMicro, nil
	case "stampnano":
		return StampNano, nil
	case "datetime":
		return DateTime, nil
	case "dateonly":
		return DateOnly, nil
	case "timeonly":
		return TimeOnly, nil
	case "unix", "unixsecond", "second", "s":
		return Unix, nil
	case "unixmilli", "unixmillisecond", "millisecond", "ms":
		return UnixMilli, nil
	case "unixmicro", "unixmicrosecond", "microsecond", "us":
		return UnixMicro, nil
	default:
		return -1, fmt.Errorf("invalid format: %s", formatStr)
	}
}

func (f Format) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, f.String())), nil
}

func (f *Format) UnmarshalJSON(data []byte) error {
	var formatStr string
	if err := json.Unmarshal(data, &formatStr); err != nil {
		return err
	}
	format, err := ParseFormat(formatStr)
	if err != nil {
		return err
	}
	*f = format
	return nil
}

func (f Format) MarshalText() ([]byte, error) {
	return []byte(f.String()), nil
}

func (f *Format) UnmarshalText(text []byte) error {
	format, err := ParseFormat(string(text))
	if err != nil {
		return err
	}
	*f = format
	return nil
}

func FormatTime(f Format, t time.Time) string {
	switch f {
	case Layout:
		return t.Format(time.Layout)
	case ANSIC:
		return t.Format(time.ANSIC)
	case UnixDate:
		return t.Format(time.UnixDate)
	case RubyDate:
		return t.Format(time.RubyDate)
	case RFC822:
		return t.Format(time.RFC822)
	case RFC822Z:
		return t.Format(time.RFC822Z)
	case RFC850:
		return t.Format(time.RFC850)
	case RFC1123:
		return t.Format(time.RFC1123)
	case RFC1123Z:
		return t.Format(time.RFC1123Z)
	case RFC3339:
		return t.Format(time.RFC3339)
	case RFC3339Nano:
		return t.Format(time.RFC3339Nano)
	case Kitchen:
		return t.Format(time.Kitchen)
	case Stamp:
		return t.Format(time.Stamp)
	case StampMilli:
		return t.Format(time.StampMilli)
	case StampMicro:
		return t.Format(time.StampMicro)
	case StampNano:
		return t.Format(time.StampNano)
	case DateTime:
		return t.Format(time.DateTime)
	case DateOnly:
		return t.Format(time.DateOnly)
	case TimeOnly:
		return t.Format(time.TimeOnly)
	case UnixMilli:
		return strconv.FormatInt(t.UnixMilli(), 10)
	case UnixMicro:
		return strconv.FormatInt(t.UnixMicro(), 10)
	case Unix:
		return strconv.FormatInt(t.Unix(), 10)
	default:
		return ""
	}
}
