package main

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Create a single instance of validator to reuse
var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

// getValidationErrorMsg returns a user-friendly error message based on the validation tag
func GetValidationErrorMsg(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		if err.Type().Kind() == reflect.String {
			return "Must be at least " + err.Param() + " characters long"
		}
		return "Must be at least " + err.Param()
	case "max":
		if err.Type().Kind() == reflect.String {
			return "Must not exceed " + err.Param() + " characters"
		}
		return "Must not exceed " + err.Param()
	case "url":
		return "Invalid URL format"
	case "e164":
		return "Invalid phone number format"
	case "containsany":
		return "Must contain at least one special character (!@#$%^&*)"
	default:
		return "Invalid value"
	}
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_578 // 1mb
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}

func sendError(w http.ResponseWriter, status int, errors []ValidationError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"errors":  errors,
	})
}

func writeJSONError(w http.ResponseWriter, status int, message string) error {
	type envelope struct {
		Error string `json:"error"`
	}

	return writeJSON(w, status, &envelope{Error: message})
}

func (app *application) validationErrorFormatter(err error) []ValidationError {
	var validationErrors []ValidationError

	for _, err := range err.(validator.ValidationErrors) {
		// Convert each validation error into our custom format
		validationErrors = append(validationErrors, ValidationError{
			Field: strings.ToLower(err.Field()),
			Error: GetValidationErrorMsg(err),
		})
	}

	return validationErrors
}
