/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
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
		Use:          "tz",
		Short:        "a timestamp offset CLI tool",
		Long:         `A CLI tool for converting to other timezones. Its mostly a wrapper around Go's time package.`,
		ValidArgs:    []string{"time"},
		Args:         cobra.ArbitraryArgs,
		Example:      "tz [FLAG]... timestamp",
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
			parsedTime, err := chrono.Parse(format, inputTime)
			if err != nil {
				return fmt.Errorf("failed to parse input time: %w", err)
			}
			outFormat, err := getOutputFormat(cmd)
			if err != nil {
				return fmt.Errorf("failed to get output-format: %w", err)
			}
			changedTime, err := ChangeOutputTimezone(cmd, parsedTime)
			if err != nil {
				return fmt.Errorf("failed to change output timezone: %w", err)
			}
			formattedTime := chrono.FormatTime(outFormat, changedTime)
			if n, err := fmt.Println(formattedTime); err != nil {
				return fmt.Errorf("failed to print parsed time after writing %d bytes: %w", n, err)
			}
			return nil
		},
	}
	//rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output.")
	defaultFormat := chrono.RFC3339.String()
	//TODO: add support for custom formats.
	rootCmd.Flags().StringP("input-format", "i", defaultFormat, "The format the time will be parsed as. (e.g. unix, rfc3339, all goland time package layouts, etc.)")
	rootCmd.Flags().StringP("output-format", "o", defaultFormat, "The format time will be output as.")
	location := time.Now().Location()
	rootCmd.Flags().StringP("output-location", "l", location.String(), "The timezone location to use (e.g. America/New_York or offsets: -5, +12:45). See https://data.iana.org/time-zones/tzdb-2021a/zone1970.tab for a list of valid locations.")
	return rootCmd
}

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
		return -1, fmt.Errorf("--input-format flag is required: %w", err)
	}
	formatStr = strings.TrimSpace(formatStr)
	if formatStr == "" {
		return -1, fmt.Errorf("--input-format flag argument cannot be empty")
	}
	f, err := chrono.ParseFormat(formatStr)
	if err != nil {
		return -1, fmt.Errorf("invalid format: %w", err)
	}
	return f, nil
}

func getOutputFormat(cmd *cobra.Command) (chrono.Format, error) {
	//if !cmd.Flags().Changed("output-format") {
	//	return getInputFormat(cmd)
	//}
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
