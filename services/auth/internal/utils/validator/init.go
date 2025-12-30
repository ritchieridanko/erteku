package validator

import "fmt"

type Validator struct{}

func Init() *Validator {
	return &Validator{}
}

func (v *Validator) Email(value *string) (bool, string) {
	if value == nil {
		return false, "Email is not provided"
	}
	if !rgxEmail.MatchString(*value) {
		return false, fmt.Sprintf("Email is invalid: %s", *value)
	}
	return true, ""
}

func (v *Validator) Password(value *string) (bool, string) {
	if value == nil {
		return false, "Password is not provided"
	}
	if *value == "" {
		return false, "Password is empty"
	}
	if len(*value) < minPasswordLength {
		return false, fmt.Sprintf("Password must be at least %d characters", minPasswordLength)
	}
	if len(*value) > maxPasswordLength {
		return false, fmt.Sprintf("Password must not exceed %d characters", maxPasswordLength)
	}
	if !rgxLowercase.MatchString(*value) {
		return false, "Password must include at least one lowercase letter"
	}
	if !rgxUppercase.MatchString(*value) {
		return false, "Password must include at least one uppercase letter"
	}
	if !rgxNumber.MatchString(*value) {
		return false, "Password must include at least one number"
	}
	if !rgxSpecialChars.MatchString(*value) {
		return false, fmt.Sprintf(
			"Password must include at least one special character: %s",
			specialChars,
		)
	}
	return true, ""
}
