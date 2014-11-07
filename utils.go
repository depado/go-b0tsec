package main

import (
	"fmt"
)

// This function takes an array of string and returns a formatted string
func GenerateTargetString(targets []string) string {
	targetsFormatted := ""
	if len(targets) > 1 {
		for i, field := range targets {
			if i == len(targets)-1 {
				targetsFormatted += fmt.Sprintf("and %v", field)
			} else if i == len(targets)-2 {
				targetsFormatted += fmt.Sprintf("%v ", field)
			} else {
				targetsFormatted += fmt.Sprintf("%v, ", field)
			}
		}
	} else {
		targetsFormatted = targets[0]
	}
	return targetsFormatted
}
