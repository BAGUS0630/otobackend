package utils

import (
	"regexp"
	"unicode"
)

// ValidateEmail mengecek format email
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidateUsername mengecek validitas username (3-20 karakter, alphanumeric + underscore)
func ValidateUsername(username string) bool {
	if len(username) < 3 || len(username) > 20 {
		return false
	}
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return usernameRegex.MatchString(username)
}

// ValidatePassword mengecek kekuatan password
// Minimal 8 karakter, harus ada uppercase, lowercase, number, dan special character
func ValidatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	hasUpper := false
	hasLower := false
	hasNumber := false
	hasSpecial := false

	for _, r := range password {
		if unicode.IsUpper(r) {
			hasUpper = true
		} else if unicode.IsLower(r) {
			hasLower = true
		} else if unicode.IsDigit(r) {
			hasNumber = true
		} else {
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

// ValidationError berisi field dan message error validasi
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidateRequired mengecek apakah string kosong
func ValidateRequired(value string, fieldName string) *ValidationError {
	if value == "" {
		return &ValidationError{
			Field:   fieldName,
			Message: fieldName + " wajib diisi",
		}
	}
	return nil
}

// ValidateMinLength mengecek panjang minimum string
func ValidateMinLength(value string, minLength int, fieldName string) *ValidationError {
	if len(value) < minLength {
		return &ValidationError{
			Field:   fieldName,
			Message: fieldName + " minimal harus " + string(rune(minLength)) + " karakter",
		}
	}
	return nil
}

// ValidateMaxLength mengecek panjang maksimum string
func ValidateMaxLength(value string, maxLength int, fieldName string) *ValidationError {
	if len(value) > maxLength {
		return &ValidationError{
			Field:   fieldName,
			Message: fieldName + " maksimal " + string(rune(maxLength)) + " karakter",
		}
	}
	return nil
}
