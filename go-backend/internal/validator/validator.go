// Package validator provides input validation functions.
package validator

import (
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

var validStatuses = map[string]bool{
	"pending":     true,
	"in-progress": true,
	"completed":   true,
}

// Email checks if the given email has a valid format.
func Email(email string) bool {
	return emailRegex.MatchString(email)
}

// Status checks if the given status is one of the allowed values.
func Status(status string) bool {
	return validStatuses[status]
}

// NonEmpty checks if a string is non-empty after trimming whitespace.
func NonEmpty(s string) bool {
	return strings.TrimSpace(s) != ""
}
