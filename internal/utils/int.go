package utils

import "math/rand"

func RandRange(min, max int) int {
	return rand.Intn(max-min) + min
}
