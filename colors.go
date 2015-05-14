package main

import (
	"github.com/fatih/color"
)

var (
	bold    = color.New(color.FgWhite).Add(color.Bold).SprintfFunc()
	red     = color.New(color.FgRed).SprintfFunc()
	redBold = color.New(color.FgRed).Add(color.Bold).SprintfFunc()
	green   = color.New(color.FgGreen).Add(color.Bold).SprintfFunc()
	cyan    = color.New(color.FgCyan).SprintfFunc()
	yellow  = color.New(color.FgYellow).Add(color.Bold).SprintfFunc()
)
