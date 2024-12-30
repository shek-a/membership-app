package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidDate(t *testing.T) {
	testCases := []struct {
		dateStr  string
		expected bool
	}{
		{"2023-10-15", true},
		{"2023-02-29", false}, // 2023 is not a leap year
		{"2020-02-29", true},  // 2020 is a leap year
		{"15-10-2023", false},
		{"2023/10/15", false},
		{"", false},
		{"2023-13-01", false}, // Invalid month
		{"2023-00-01", false}, // Invalid month
		{"2023-12-32", false}, // Invalid day
	}

	for _, tc := range testCases {
		t.Run(tc.dateStr, func(t *testing.T) {
			result := IsValidDate(tc.dateStr)
			assert.Equal(t, tc.expected, result)
		})
	}
}
