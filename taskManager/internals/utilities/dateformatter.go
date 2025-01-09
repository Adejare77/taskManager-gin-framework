package utilities

import (
	"fmt"
	"strings"
	"time"
)

func TimeManipulator(input string) string {
	layout := "2006-01-02T15:04" // Constant layout

	if strings.Contains(input, "day") || strings.Contains(input, "hour") || strings.Contains(input, "minute") {
		var number int
		var unit string

		fmt.Sscanf(input, "%d %s", &number, &unit)

		currentTime := time.Now()

		switch {
		case strings.Contains(input, "day"):
			return (currentTime.Add(time.Duration(number) * 24 * time.Hour)).Format(layout)
		case strings.Contains(input, "minutes"):
			return (currentTime.Add(time.Duration(number) * time.Minute)).Format(layout)
		case strings.Contains(input, "hour"):
			return (currentTime.Add(time.Duration(number) * time.Hour)).Format(layout)
		}
	}

	parsedTime, _ := time.Parse(layout, input)
	return parsedTime.Format(layout)

}
