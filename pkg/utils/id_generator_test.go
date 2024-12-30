package utils

import (
	"testing"
)

func TestGenerateRandomNumber(t *testing.T) {
	for i := 0; i < 100; i++ { // Run the test multiple times to ensure randomness
		num := GenerateRandomNumber()
		if num < 111111 || num > 999999 {
			t.Errorf("Generated number %d is out of range", num)
		}
	}
}
