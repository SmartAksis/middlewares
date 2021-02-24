package utils

import "regexp"

func IsEmail(email string) bool {
	regexEmailPattern := `^[a-z0-9.]+@[a-z0-9]+\.[a-z]+(\.[a-z]+)?$`
	return regexp.MustCompile(regexEmailPattern).MatchString(email)
}