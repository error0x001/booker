package utils

import "regexp"

var (
	emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
)

func IsEmailValid(e string) bool {
	return emailRegex.MatchString(e)
}
