package utilities

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func startDateValidation(fl validator.FieldLevel) bool {
	startDate := fl.Field().String()

	layout := "2006-01-02 15:04" // Constant layout

	_, err := time.Parse(layout, startDate)
	if err != nil {
		return err == nil
	}
	return true
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

	layout := "2006-01-02 15:04" // Constant layout

	_, err := time.Parse(layout, date)
	if err != nil {
		return err == nil
	}

	return true
}

func RegisterValidation() {
	// Register the above Validations
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("startdate", startDateValidation)
		v.RegisterValidation("duedate", dueDateValidation)
	}
}

// func ValidationError(err validator.ValidationErrors) []any {
// 	var errorDetails []any
// 	for _, fieldError := range err {
// 		if fieldError.Tag() == "required" {
// 			errorDetails = append(errorDetails, fmt.Sprintf(
// 				"missing %s field", fieldError.Field(),
// 			))
// 		} else if fieldError.Tag() == "duedate" {
// 			errorDetails = append(errorDetails,
// 				map[string]any{
// 					"dueDate format": "`YYYY-MM-DD HH:MM` e.g., 2024-05-19 22:15," +
// 						"`x day(s)` e.g., 3 days (relative to the current time)," +
// 						"`x hour(s)` e.g., 5 hours (relative to the current time)",
// 				},
// 			)
// 		} else if fieldError.Tag() == "startdate" {
// 			errorDetails = append(errorDetails,
// 				"startDate format: `YYYY-MM-DD HH:MM` e.g., 2024-05-19 22:15",
// 			)
// 		} else if fieldError.Tag() == "email" {
// 			errorDetails = append(errorDetails, "Invalid email Format")
// 		} else {
// 			errorDetails = append(errorDetails, err.Error())
// 		}
// 	}
// 	return errorDetails
// }

func ValidationError(err error) interface{} {
	if err == nil {
		return nil
	}

	// Handle Validation errors
	if fieldErrors, ok := err.(validator.ValidationErrors); ok {
		errorMessages := make(map[string]string)
		for _, fieldError := range fieldErrors {
			field := strings.ToLower(fieldError.Field())
			errorMessages[field] = fmt.Sprintf("Invalid value for %s: %s", field, fieldError.Tag())
		}
		return errorMessages
	}

	// Handle other types of errors (e.g JSON parsing errors)
	return err.Error()
}
