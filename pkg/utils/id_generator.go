package utils

import (
	"math/rand"
	"time"
)

func GenerateRandomNumber() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(999999-111111+1) + 111111
}
