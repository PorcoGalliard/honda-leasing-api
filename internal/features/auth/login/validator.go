package login

import (
	"fmt"
	"regexp"
)

func ValidateLoginRequest(req LoginRequest) error {
	// Validate phone number format
	phoneRegex := regexp.MustCompile(`^(\+62|62|0)[0-9]{9,12}$`)
	if !phoneRegex.MatchString(req.PhoneNumber) {
		return fmt.Errorf("invalid phone number format")
	}

	// Validate PIN format
	pinRegex := regexp.MustCompile(`^[0-9]{6}$`)
	if !pinRegex.MatchString(req.PIN) {
		return fmt.Errorf("PIN must be 6 digits")
	}

	return nil
}
