package utils

import (
	"errors"
	"regexp"
	"strings"
)

// PasswordPolicy holds password requirements
type PasswordPolicy struct {
	MinLength        int
	RequireUppercase bool
	RequireLowercase bool
	RequireNumbers   bool
	RequireSpecial   bool
	ForbiddenWords   []string
}

// DefaultPasswordPolicy returns the default password policy
func DefaultPasswordPolicy() PasswordPolicy {
	return PasswordPolicy{
		MinLength:        12,
		RequireUppercase: true,
		RequireLowercase: true,
		RequireNumbers:   true,
		RequireSpecial:   true,
		ForbiddenWords: []string{
			"password", "admin", "123456", "qwerty", "abc123",
			"password123", "admin123", "user", "login", "welcome",
		},
	}
}

// ValidatePassword validates a password against the security policy
func ValidatePassword(password string) error {
	policy := DefaultPasswordPolicy()

	// Check minimum length
	if len(password) < policy.MinLength {
		return errors.New("password must be at least 12 characters long")
	}

	// Check for uppercase letters
	if policy.RequireUppercase {
		if matched, _ := regexp.MatchString(`[A-Z]`, password); !matched {
			return errors.New("password must contain at least one uppercase letter")
		}
	}

	// Check for lowercase letters
	if policy.RequireLowercase {
		if matched, _ := regexp.MatchString(`[a-z]`, password); !matched {
			return errors.New("password must contain at least one lowercase letter")
		}
	}

	// Check for numbers
	if policy.RequireNumbers {
		if matched, _ := regexp.MatchString(`[0-9]`, password); !matched {
			return errors.New("password must contain at least one number")
		}
	}

	// Check for special characters
	if policy.RequireSpecial {
		if matched, _ := regexp.MatchString(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`, password); !matched {
			return errors.New("password must contain at least one special character (!@#$%^&*)")
		}
	}

	// Check for forbidden words
	passwordLower := strings.ToLower(password)
	for _, word := range policy.ForbiddenWords {
		if strings.Contains(passwordLower, strings.ToLower(word)) {
			return errors.New("password contains forbidden words")
		}
	}

	// Check for common patterns
	if matched, _ := regexp.MatchString(`(.)\1{2,}`, password); matched {
		return errors.New("password cannot contain more than 2 consecutive identical characters")
	}

	// Check for sequential characters
	if containsSequentialChars(password) {
		return errors.New("password cannot contain sequential characters (abc, 123, etc.)")
	}

	return nil
}

// containsSequentialChars checks for sequential character patterns
func containsSequentialChars(password string) bool {
	sequences := []string{
		"abcdefghijklmnopqrstuvwxyz",
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		"0123456789",
		"qwertyuiop",
		"asdfghjkl",
		"zxcvbnm",
	}

	for _, seq := range sequences {
		for i := 0; i <= len(seq)-3; i++ {
			if strings.Contains(strings.ToLower(password), seq[i:i+3]) {
				return true
			}
		}
	}

	return false
}

// PasswordStrength calculates password strength score (0-100)
func PasswordStrength(password string) int {
	score := 0

	// Length bonus
	if len(password) >= 8 {
		score += 10
	}
	if len(password) >= 12 {
		score += 15
	}
	if len(password) >= 16 {
		score += 10
	}

	// Character variety bonus
	if matched, _ := regexp.MatchString(`[a-z]`, password); matched {
		score += 10
	}
	if matched, _ := regexp.MatchString(`[A-Z]`, password); matched {
		score += 10
	}
	if matched, _ := regexp.MatchString(`[0-9]`, password); matched {
		score += 10
	}
	if matched, _ := regexp.MatchString(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`, password); matched {
		score += 15
	}

	// Complexity bonus
	uniqueChars := make(map[rune]bool)
	for _, char := range password {
		uniqueChars[char] = true
	}
	if len(uniqueChars) > len(password)/2 {
		score += 10
	}

	// Penalty for common patterns
	if containsSequentialChars(password) {
		score -= 20
	}
	if matched, _ := regexp.MatchString(`(.)\1{2,}`, password); matched {
		score -= 15
	}

	// Ensure score is between 0-100
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return score
}
