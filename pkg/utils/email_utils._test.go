package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidEmail(t *testing.T) {
	testCases := []struct {
		emailStr string
		expected bool
	}{
		{"test@example.com", true},
		{"user.name+tag+sorting@example.com", true},
		{"user.name@example.co.uk", true},
		{"invalid-email", false},
		{"@example.com", false},
		{"test@.com", false},
		{"test@domain", false},
		{"test@domain.c", false},
		{"test@domain.corporate", true},
	}

	for _, tc := range testCases {
		t.Run(tc.emailStr, func(t *testing.T) {
			result := IsValidEmail(tc.emailStr)
			assert.Equal(t, tc.expected, result)
		})
	}
}
