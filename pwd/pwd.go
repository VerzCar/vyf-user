package pwd

import (
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

const MIN_PWD_LENGTH = 8

var regexpLowerCase = regexp.MustCompile(`[a-z]{1}`)
var regexpUpperCase = regexp.MustCompile(`[A-Z]{1}`)
var regexpNumberCase = regexp.MustCompile(`[\d]{1}`)
var regexpSpecialCase = regexp.MustCompile(`[\W+]{1}`)

// Hash from the given plain text password.
func Hash(pwd string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(pwd), 10)

	return string(hashed), err
}

// Verify verifies, that the given plain text password and hashed password matches.
// If they match nil will be returned, otherwise an error.
func Verify(pwd string, hashed string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(pwd))
}

// Complexity checks the required complexity of the
// given password.
// Returns true if the desired complexity match the given value.
func Complexity(pwd string) bool {
	isComplexityGiven := true

	if len(pwd) < MIN_PWD_LENGTH {
		isComplexityGiven = false
		return isComplexityGiven
	}

	if !regexpLowerCase.MatchString(pwd) {
		isComplexityGiven = false
	}

	if !regexpUpperCase.MatchString(pwd) {
		isComplexityGiven = false
	}

	if !regexpNumberCase.MatchString(pwd) {
		if !regexpSpecialCase.MatchString(pwd) {
			isComplexityGiven = false
		}
	}

	return isComplexityGiven
}
