package utilities

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func statusValidation(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	return status == "pending" || status == "in-progress"
}

func dueDateValidation(fl validator.FieldLevel) bool {
	date := fl.Field().String()
	if strings.Contains(date, "day") || strings.Contains(date, "hour") || strings.Contains(date, "minute") {
		var number int
		var unit string

		_, err := fmt.Sscanf(date, "%d %s", &number, &unit)
		if err != nil {
			return false
		}

		switch {
		case strings.Contains(date, "day"):
			return true
		case strings.Contains(date, "minute"):
			return true
		case strings.Contains(date, "hour"):
			return true
		default:
			return false
		}
	}

	layout := "2006-01-02T15:04" // Constant layout

	parsedTime, err := time.Parse(layout, date)
	if err == nil {
		return parsedTime.Format(layout) > time.Now().String()
	}
	return false
}

func RegisterValidation() {
	// Register the above Validations
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("status", statusValidation)
		v.RegisterValidation("dueDate", dueDateValidation)
	}
}

func ValidationError(err validator.ValidationErrors) []string {
	var errorDetails []string
	for _, fieldError := range err {
		if fieldError.Tag() == "required" {
			errorDetails = append(errorDetails, fmt.Sprintf(
				"missing %s field", fieldError.Field(),
			))
		} else if fieldError.Tag() == "status" {
			errorDetails = append(errorDetails,
				"status can only be `pending` or `in-progress`",
			)
		} else if fieldError.Tag() == "dueDate" {
			errorDetails = append(errorDetails,
				"dueDate format:\n`YYYY-MM-DDTHH:MM` e.g., 2024-05-19T22:15,\n`x day(s)` e.g., 3 days (relative to the current time)\n`x hour(s)` e.g., 5 hours (relative to the current time)",
			)
		}
	}
	return errorDetails
}
