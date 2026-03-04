package performance

import (
	"strings"
)

// OptimizedValidator provides validation with pre-compiled regex patterns
type OptimizedValidator struct {
	regexCache *RegexCache
}

// NewOptimizedValidator creates a new OptimizedValidator instance
func NewOptimizedValidator() *OptimizedValidator {
	validator := &OptimizedValidator{
		regexCache: NewRegexCache(),
	}

	// Pre-compile common patterns
	commonPatterns := []string{
		`^1[3-9]\d{9}$`, // Chinese phone
		`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, // Email
		`[a-zA-Z]`,               // Has letter
		`[0-9]`,                  // Has number
		`[A-Z]`,                  // Has uppercase
		`[a-z]`,                  // Has lowercase
		`[!@#$%^&*(),.?":{}|<>]`, // Has special char
		`[<>""'\\]`,              // Dangerous chars
	}

	validator.regexCache.Precompile(commonPatterns)
	return validator
}

// ValidatePhone validates Chinese phone numbers using pre-compiled regex
func (ov *OptimizedValidator) ValidatePhone(phone string) bool {
	if len(phone) != 11 {
		return false
	}

	regex := ov.regexCache.MustGet(`^1[3-9]\d{9}$`)
	return regex.MatchString(phone)
}

// ValidateEmail validates email addresses using pre-compiled regex
func (ov *OptimizedValidator) ValidateEmail(email string) bool {
	regex := ov.regexCache.MustGet(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return regex.MatchString(email)
}

// ValidatePasswordStrength validates password strength using pre-compiled regex
func (ov *OptimizedValidator) ValidatePasswordStrength(password string) (bool, string) {
	if len(password) < 6 {
		return false, "密码长度不能小于6位"
	}

	if len(password) > 32 {
		return false, "密码长度不能大于32位"
	}

	hasUpper := ov.regexCache.MustGet(`[A-Z]`).MatchString(password)
	hasLower := ov.regexCache.MustGet(`[a-z]`).MatchString(password)
	hasNumber := ov.regexCache.MustGet(`[0-9]`).MatchString(password)
	hasSpecial := ov.regexCache.MustGet(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)

	strength := 0
	if hasUpper {
		strength++
	}
	if hasLower {
		strength++
	}
	if hasNumber {
		strength++
	}
	if hasSpecial {
		strength++
	}

	if strength < 2 {
		return false, "密码必须包含字母和数字，建议包含特殊字符"
	}

	return true, ""
}

// SanitizeInput cleans input data using pre-compiled regex
func (ov *OptimizedValidator) SanitizeInput(input string) string {
	// Remove leading/trailing spaces
	input = strings.TrimSpace(input)

	// Remove dangerous characters using pre-compiled regex
	regex := ov.regexCache.MustGet(`[<>""'\\]`)
	input = regex.ReplaceAllString(input, "")

	return input
}

// ValidateEmployeeNo validates employee number format
func (ov *OptimizedValidator) ValidateEmployeeNo(no string) bool {
	// Employee number must be 6 digits
	if len(no) != 6 {
		return false
	}

	// Check if all characters are digits
	for _, ch := range no {
		if ch < '0' || ch > '9' {
			return false
		}
	}

	return true
}

// Global instance for convenience
var globalValidator = NewOptimizedValidator()

// GetOptimizedValidator returns the global OptimizedValidator instance
func GetOptimizedValidator() *OptimizedValidator {
	return globalValidator
}
