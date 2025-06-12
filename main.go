/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"os"

	"github.com/jakobii/tz/internal/cmd"
	"github.com/jakobii/tz/internal/parsers"
)

func main() {
	tp := parsers.NewTimeParser(parsers.DefaultLayouts)
	if err := cmd.NewRootCommand(tp).Execute(); err != nil {
		os.Exit(1)
	}
}
