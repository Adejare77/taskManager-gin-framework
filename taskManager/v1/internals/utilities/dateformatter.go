package utilities

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

func StartTimeManipulator(input string) (time.Time, error) {
	if input == "" {
		return time.Now(), nil
	}

	layout := "2006-01-02 15:04" // Constant layout
	parsedTime, err := time.Parse(layout, input)
	if err != nil {
		fmt.Println("DATEFORMATTER", input) // Meant to check where the error is
		fmt.Println("Error: ", err)         // Takes only the required format
		return time.Time{}, err
	}

	if parsedTime.Before(time.Now()) {
		return time.Time{}, errors.New("start date cannot be in the past")
	}

	return parsedTime, nil
}

func DueTimeManipulator(input string, startDate time.Time) (time.Time, error) {
	layout := "2006-01-02 15:04" // Constant layout

	if strings.Contains(input, "day") || strings.Contains(input, "hour") || strings.Contains(input, "minute") {
		var number int
		var unit string

		fmt.Sscanf(input, "%d %s", &number, &unit)

		switch {
		case strings.Contains(input, "day"):
			return startDate.Add(time.Duration(number) * 24 * time.Hour), nil
		case strings.Contains(input, "minutes"):
			return startDate.Add(time.Duration(number) * time.Minute), nil
		case strings.Contains(input, "hour"):
			return startDate.Add(time.Duration(number) * time.Hour), nil
		}
	}

	parsedTime, _ := time.Parse(layout, input)

	if parsedTime.Before(startDate) {
		return time.Time{}, errors.New("due date must be After start date")
	}

	return parsedTime, nil
}
