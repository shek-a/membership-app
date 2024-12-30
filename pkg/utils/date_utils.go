package utils

import (
	"time"
)

const dateFormat = "2006-01-02"

func IsValidDate(dateStr string) bool {
	_, err := time.Parse(dateFormat, dateStr)
	return err == nil
}
