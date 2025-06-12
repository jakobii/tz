/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"os"

	"github.com/jakobii/tz/internal/cmd"
)

func main() {
	root := cmd.NewRootCommand()
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
