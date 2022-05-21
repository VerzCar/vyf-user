package pwd

import (
	assertPkg "github.com/stretchr/testify/assert"
	"testing"
)

func TestHashPassword(t *testing.T) {
	assert := assertPkg.New(t)

	password := "Commando1"

	hashed, err := Hash(password)

	if err != nil {
		t.Fatalf(`Could not hash password: %s`, err)
	}

	assert.NotEmpty(hashed)
}

func TestVerifyPassword(t *testing.T) {
	assert := assertPkg.New(t)

	password := "Commando1"

	hashed, err := Hash(password)

	if err != nil {
		t.Fatalf(`Could not hash password: %s`, err)
	}

	// password must match
	assert.Nil(Verify(password, hashed))

	invalidPassword := "JohnDoe23"

	// password must not match
	assert.NotNil(Verify(invalidPassword, hashed))
}

func TestCheckPasswordComplexity(t *testing.T) {
	assert := assertPkg.New(t)

	validPassword := "Commando1"

	// complexity matches
	assert.Equal(true, Complexity(validPassword))

	validPassword = "Commando#"

	// complexity matches
	assert.Equal(true, Complexity(validPassword))

	validPassword = "Commando1#"

	// complexity matches
	assert.Equal(true, Complexity(validPassword))

	validPassword = "Command1"

	// complexity matches
	assert.Equal(true, Complexity(validPassword))

	validPassword = "Command#"

	// complexity matches
	assert.Equal(true, Complexity(validPassword))

	invalidPassword := "comm"

	// complexity does not match
	assert.Equal(false, Complexity(invalidPassword))

	invalidPassword = "Comm"

	// complexity does not match
	assert.Equal(false, Complexity(invalidPassword))

	invalidPassword = "Commando"

	// complexity does not match
	assert.Equal(false, Complexity(invalidPassword))

	invalidPassword = "commando1"

	// complexity does not match
	assert.Equal(false, Complexity(invalidPassword))

	invalidPassword = "commando#"

	// complexity does not match
	assert.Equal(false, Complexity(invalidPassword))

	invalidPassword = "Comman#"

	// complexity does not match
	assert.Equal(false, Complexity(invalidPassword))

	invalidPassword = "commando#1"

	// complexity does not match
	assert.Equal(false, Complexity(invalidPassword))

}
