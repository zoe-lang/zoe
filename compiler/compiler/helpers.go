package zoe

import "github.com/fatih/color"

var fred = color.New(color.FgRed, color.Bold)
var red = fred.SprintFunc()
var fgreen = color.New(color.FgGreen)
var green = fgreen.SprintFunc()
var cyan = color.New(color.FgCyan).SprintFunc()
var yel = color.New(color.FgYellow).SprintFunc()
var mag = color.New(color.FgMagenta).SprintFunc()
var blue = color.New(color.FgBlue).SprintFunc()
var grey = color.New(color.Faint).SprintFunc()

func maxInt(a, b uint32) uint32 {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}
