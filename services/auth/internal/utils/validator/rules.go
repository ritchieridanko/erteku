package validator

import "regexp"

var (
	rgxEmail        = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	rgxLowercase    = regexp.MustCompile(`[a-z]`)
	rgxNumber       = regexp.MustCompile(`[0-9]`)
	rgxSpecialChars = regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`)
	rgxUppercase    = regexp.MustCompile(`[A-Z]`)
)

const (
	minPasswordLength int = 8
	maxPasswordLength int = 50

	specialChars string = `!@#$%^&*(),.?":{}|<>`
)
