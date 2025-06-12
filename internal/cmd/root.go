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

	"github.com/jakobii/tz/internal/parsers"
	"github.com/spf13/cobra"
)

type timeParser interface {
	ParseTime(v string) (*parsers.TimeResult, error)
}

func NewRootCommand(parser timeParser) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:       "tz",
		Short:     "a timezone CLI tool",
		Long:      `A CLI tool for converting to other timezones. Its mostly a wrapper around Go's time package.`,
		ValidArgs: []string{"time"},
		Args:      cobra.ArbitraryArgs,
		Example:   "tz --time 2025-06-12T02:57:30Z",
		RunE: func(cmd *cobra.Command, args []string) error {
			isVerbose := cmd.PersistentFlags().Changed("verbose")
			var inputTime string
			if hasStdin() {
				if b, err := io.ReadAll(io.LimitReader(os.Stdin, 1024)); err == nil {
					inputTime = string(b)
				}
			} else {
				_time, err := cmd.Flags().GetString("time")
				if err != nil {
					return fmt.Errorf("--time flag is required: %w", err)
				}
				inputTime = _time
			}
			inputTime = strings.TrimSpace(inputTime)
			if inputTime == "" {
				return fmt.Errorf("--time flag argument cannot be empty")
			}
			result, err := parser.ParseTime(inputTime)
			if err != nil {
				return fmt.Errorf("failed to parse time: %w", err)
			}
			if isVerbose {
				if n, err := fmt.Fprintf(os.Stdout, "Input time: %s\n", inputTime); err != nil {
					return fmt.Errorf("failed to print input time after writing %d bytes: %w", n, err)
				}
				if n, err := fmt.Fprintf(os.Stdout, "Tried as: %+v\n", result.TriedAs); err != nil {
					return fmt.Errorf("failed to print tried as after writing %d bytes: %w", n, err)
				}
				if n, err := fmt.Fprintf(os.Stdout, "Parsed as: %s\n", result.ParsedAs); err != nil {
					return fmt.Errorf("failed to print parsed as after writing %d bytes: %w", n, err)
				}
				if n, err := fmt.Fprintf(os.Stdout, "Formatted as: %s\n", parsers.RFC3339); err != nil {
					return fmt.Errorf("failed to print parsed time after writing %d bytes: %w", n, err)
				}
			}

			if n, err := fmt.Fprintln(os.Stdout, result.Time.Local().Format(time.RFC3339)); err != nil {
				return fmt.Errorf("failed to print parsed time after writing %d bytes: %w", n, err)
			}
			return nil
		},
	}
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output.")
	rootCmd.Flags().StringP("time", "t", "", "The input time to parse.")
	// TODO: add support for layout and location flags
	//rootCmd.Flags().StringP("layout", "l", "", "The layout to use for parsing the input time. Defaults to Go's time package layouts https://pkg.go.dev/time#pkg-constants.")
	//rootCmd.Flags().StringP("location", "z", "", "The timezone location to use (e.g. America/New_York). See https://data.iana.org/time-zones/tzdb-2021a/zone1970.tab for a list of valid locations.")
	return rootCmd
}

func hasStdin() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (info.Mode() & os.ModeCharDevice) == 0
}
