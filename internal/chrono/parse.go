package chrono

import (
	"strconv"
	"time"
)

func Parse(f Format, v string) (time.Time, error) {
	switch f {
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
	case UnixMicro:
		vInt, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return time.Time{}, err
		}
		return time.UnixMicro(vInt), nil
	case Unix:
		vInt, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return time.Time{}, err
		}
		return time.Unix(vInt, 0), nil
	default:
		return time.Time{}, errUnknownLayout
	}
}
