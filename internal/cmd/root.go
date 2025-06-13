package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/jakobii/tz/internal/chrono"
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "tz [flags]... [Timestamp]",
		Short: "a timestamp offset CLI tool",
		Long: `A CLI tool for converting to other timezones. ` +
			`Its mostly a wrapper around Go's time package. ` +
			`Accepts timestamp from stdin or as its last argument. ` +
			`Go's time package uses a single date to communicate its various formatting options. ` +
			`See https://pkg.go.dev/time#pkg-constants for more details. ` +
			formatRules,
		Args: cobra.ArbitraryArgs,
		Example: "tz '2009-11-10T23:00:00Z'                                # change UTC to local.\n" +
			"echo '2009-11-10T23:00:00Z' | tz                         # pipe UTC to local.\n" +
			"tz -l America/New_York '2009-11-10T23:00:00Z'            # change timezone to NY.\n" +
			"tz -l -5 '2009-11-10T23:00:00Z'                          # change timezone to specific offset '-5' hours.\n" +
			"tz -l +12:45 '2009-11-10T23:00:00Z'                      # change timezone to specific offset.\n" +
			"tz -i RubyDate -o ms 'Tue Nov 10 23:00:00 +0000 2009'    # RubyDate to milliseconds since epoch.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			//isVerbose := cmd.PersistentFlags().Changed("verbose")
			inputTime, err := getTimeParam(args)
			if err != nil {
				return fmt.Errorf("failed to get input time: %w", err)
			}
			format, err := getInputFormat(cmd)
			if err != nil {
				return fmt.Errorf("failed to get format: %w", err)
			}
			var parsedTime time.Time
			if format == chrono.Unknown {
				customFormat, err := cmd.Flags().GetString("input-format")
				if err != nil {
					return fmt.Errorf("failed to get custom format: %w", err)
				}
				t, err := time.Parse(strings.TrimSpace(customFormat), inputTime)
				if err != nil {
					return fmt.Errorf("failed to parse input time with custom format: %w", err)
				}
				parsedTime = t
			} else {
				t, err := chrono.Parse(format, inputTime)
				if err != nil {
					return fmt.Errorf("failed to parse input time: %w", err)
				}
				parsedTime = t
			}
			outFormat, err := getOutputFormat(cmd)
			if err != nil {
				return fmt.Errorf("failed to get output-format: %w", err)
			}
			changedTime, err := ChangeOutputTimezone(cmd, parsedTime)
			if err != nil {
				return fmt.Errorf("failed to change output timezone: %w", err)
			}
			var formattedTimeOutput string
			if outFormat == chrono.Unknown {
				customFormat, err := cmd.Flags().GetString("output-format")
				if err != nil {
					return fmt.Errorf("failed to get custom output format: %w", err)
				}
				formattedTimeOutput = changedTime.Format(strings.TrimSpace(customFormat))
			} else {
				formattedTimeOutput = chrono.FormatTime(outFormat, changedTime)
			}
			if n, err := fmt.Println(formattedTimeOutput); err != nil {
				return fmt.Errorf("failed to print parsed time after writing %d bytes: %w", n, err)
			}
			return nil
		},
	}
	//rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output.")
	defaultFormat := chrono.RFC3339.String()
	//TODO: add support for custom formats.
	rootCmd.Flags().StringP("input-format", "i", defaultFormat, "The format the time will be parsed as. May be one of the predefined formats or a custom format string.")
	rootCmd.Flags().StringP("output-format", "o", defaultFormat, "The format time will be output as. May be one of the predefined formats or a custom format string.")
	rootCmd.Flags().StringP("output-location", "l", location.String(), "The timezone location to change to (e.g. America/New_York or offsets: -5, +12:45). See https://data.iana.org/time-zones/tzdb-2021a/zone1970.tab for a list of valid locations.")
	return rootCmd
}

const formatRules = `

	summary of the formatting components:

	Year: "2006" "06"
	Month: "Jan" "January" "01" "1"
	Day of the week: "Mon" "Monday"
	Day of the month: "2" "_2" "02"
	Day of the year: "__2" "002"
	Hour: "15" "3" "03" (PM or AM)
	Minute: "4" "04"
	Second: "5" "05"
	AM/PM mark: "PM"

	Numeric time zone offsets:

	"-0700"     ±hhmm
	"-07:00"    ±hh:mm
	"-07"       ±hh
	"-070000"   ±hhmmss
	"-07:00:00" ±hh:mm:ss

	ISO 8601 behavior:

	"Z0700"      Z or ±hhmm
	"Z07:00"     Z or ±hh:mm
	"Z07"        Z or ±hh
	"Z070000"    Z or ±hhmmss
	"Z07:00:00"  Z or ±hh:mm:ss

	Predefined formats:

	ANSIC       = "Mon Jan _2 15:04:05 2006"
	UnixDate    = "Mon Jan _2 15:04:05 MST 2006"
	RubyDate    = "Mon Jan 02 15:04:05 -0700 2006"
	RFC822      = "02 Jan 06 15:04 MST"
	RFC822Z     = "02 Jan 06 15:04 -0700" 
	RFC850      = "Monday, 02-Jan-06 15:04:05 MST"
	RFC1123     = "Mon, 02 Jan 2006 15:04:05 MST"
	RFC1123Z    = "Mon, 02 Jan 2006 15:04:05 -0700" 
	RFC3339     = "2006-01-02T15:04:05Z07:00"
	RFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"
	Kitchen     = "3:04PM"
	Stamp      = "Jan _2 15:04:05"
	StampMilli = "Jan _2 15:04:05.000"
	StampMicro = "Jan _2 15:04:05.000000"
	StampNano  = "Jan _2 15:04:05.000000000"
	DateTime   = "2006-01-02 15:04:05"
	DateOnly   = "2006-01-02"
	TimeOnly   = "15:04:05"

	Unix formats (These do not follow Go's time package layout rules):

	Unix       Number of seconds since January 1, 1970 UTC.
	UnixMilli  Number of milliseconds since January 1, 1970 UTC.
	UnixMicro  Number of microseconds since January 1, 1970 UTC.
`

func hasStdin() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (info.Mode() & os.ModeCharDevice) == 0
}

func getTimeParam(args []string) (string, error) {
	var inputTime string
	if hasStdin() {
		if b, err := io.ReadAll(io.LimitReader(os.Stdin, 1024)); err == nil {
			inputTime = string(b)
		}
	} else {
		if len(args) < 1 {
			return "", fmt.Errorf("timestamp argument is required")
		}
		timestamp := args[len(args)-1]
		inputTime = timestamp
	}
	inputTime = strings.TrimSpace(inputTime)
	if inputTime == "" {
		return "", fmt.Errorf("--time flag argument cannot be empty")
	}
	return inputTime, nil
}

func getInputFormat(cmd *cobra.Command) (chrono.Format, error) {
	formatStr, err := cmd.Flags().GetString("input-format")
	if err != nil {
		return 0, fmt.Errorf("--input-format flag is required: %w", err)
	}
	formatStr = strings.TrimSpace(formatStr)
	if formatStr == "" {
		return 0, fmt.Errorf("--input-format flag argument cannot be empty")
	}
	f, err := chrono.ParseFormat(formatStr)
	if err != nil {
		return chrono.Unknown, nil
	}
	return f, nil
}

func getOutputFormat(cmd *cobra.Command) (chrono.Format, error) {
	formatStr, err := cmd.Flags().GetString("output-format")
	if err != nil {
		return -1, fmt.Errorf("--output-format flag is required: %w", err)
	}
	formatStr = strings.TrimSpace(formatStr)
	if formatStr == "" {
		return -1, fmt.Errorf("--output-format flag argument cannot be empty")
	}
	f, err := chrono.ParseFormat(formatStr)
	if err != nil {
		return -1, fmt.Errorf("invalid format: %w", err)
	}

	return f, nil
}

func ChangeOutputTimezone(cmd *cobra.Command, t time.Time) (time.Time, error) {
	zoneStr, err := cmd.Flags().GetString("output-location")
	if err != nil {
		return time.Time{}, fmt.Errorf("--output-location flag is required: %w", err)
	}
	zoneStr = strings.TrimSpace(zoneStr)
	if zoneStr == "" {
		return time.Time{}, fmt.Errorf("--output-location flag argument cannot be empty")
	}
	var timezone *time.Location
	if seconds, err := chrono.ParseZoneOffset(zoneStr); err == nil {
		timezone = time.FixedZone("", seconds)
	} else {
		zone, err := time.LoadLocation(zoneStr)
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid timezone location: %w", err)
		}
		timezone = zone
	}
	return t.In(timezone), nil
}
