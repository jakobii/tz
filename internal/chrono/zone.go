package chrono

import (
	"regexp"
	"strconv"
)

var reZoneOffset = regexp.MustCompile(`^([+-])(\d{1,2})(?::?(\d{2}))?$`)

func ParseZoneOffset(offset string) (int, error) {
	matches := reZoneOffset.FindStringSubmatch(offset)
	if len(matches) == 0 {
		return 0, errUnknownFormat
	}

	sign := 1
	if matches[1] == "-" {
		sign = -1
	}

	hours, err := strconv.Atoi(matches[2])
	if err != nil {
		return 0, err
	}

	minutes := 0
	if len(matches) > 3 && matches[3] != "" {
		minutes, err = strconv.Atoi(matches[3])
		if err != nil {
			return 0, err
		}
	}

	return sign * (hours*3600 + minutes*60), nil
}

func IsValidZoneOffset(offset string) bool {
	_, err := ParseZoneOffset(offset)
	return err == nil
}
