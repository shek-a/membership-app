package utils

import "regexp"

func IsValidEmail(emailStr string) bool {
	// Define a regular expression for validating an email address
	const emailRegexPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	emailRegex := regexp.MustCompile(emailRegexPattern)
	return emailRegex.MatchString(emailStr)
}
