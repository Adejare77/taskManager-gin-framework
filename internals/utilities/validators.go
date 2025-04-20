package utilities

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

var layout = "2006-01-02 15:04"

func CompareDates(startDate *JSONTime, dueDate JSONTime) (JSONTime, error) {
	var startTime time.Time
	dueTime := time.Time(dueDate)
	currentTime := time.Now()

	if startDate == nil {
		startTime = currentTime

	} else {
		startTime = time.Time(*startDate)
		if startTime.Before(currentTime) {
			return JSONTime{}, fmt.Errorf("start_date (%s) cannot be in past of current time (%s)", startTime, currentTime)
		}
	}

	if dueTime.Before(startTime) {
		return JSONTime{}, fmt.Errorf("due_date (%s) cannot be in past of start_date (%s)", dueTime, startTime)
	}

	return JSONTime(startTime), nil
}

func ValidationError(err error) []string {
	var errorDetails []string
	fieldErrors, _ := err.(validator.ValidationErrors)
	validationErr := "validation failed:"

	for _, fieldError := range fieldErrors {
		field := strings.ToLower(fieldError.Field())
		if fieldError.Tag() == "required" {
			errorDetails = append(errorDetails, fmt.Sprintf("%s missing %s field", validationErr, field))
		} else if fieldError.Tag() == "numeric" {
			errorDetails = append(errorDetails, fmt.Sprintf("%s %s can only be a numeric value", validationErr, field))
		} else if fieldError.Tag() == "min" {
			errorDetails = append(errorDetails, fmt.Sprintf("%s %s cannot be less than 1", validationErr, field))
		} else if fieldError.Tag() == "oneof" {
			errorDetails = append(errorDetails, fmt.Sprintf("%s `%s` can only be one of this fields: pending, in-progress, completed", validationErr, field))
		} else if fieldError.Tag() == "uuid" {
			errorDetails = append(errorDetails, fmt.Sprintf("%s invalid post uuid", validationErr))
		} else {
			errorDetails = append(errorDetails, fmt.Sprintf("%s unresolved validation", validationErr))
		}
	}
	return errorDetails

}
