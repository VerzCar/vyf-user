package validate

import (
	assertPkg "github.com/stretchr/testify/assert"
	"testing"
)

type phoneNumber struct {
	Number      string
	CountryCode string
}

func TestValidatePhoneNumber(t *testing.T) {
	assert := assertPkg.New(t)

	validPhoneNumber := phoneNumber{
		Number:      "017683464774",
		CountryCode: "DE",
	}

	formattedNumber, err := PhoneNumber(validPhoneNumber.Number, validPhoneNumber.CountryCode)

	assert.Equal("+4917683464774", formattedNumber)
	assert.Nil(err)

	invalidPhoneNumber := phoneNumber{
		Number:      "255834643423323174",
		CountryCode: "DE",
	}

	formattedNumber, err = PhoneNumber(invalidPhoneNumber.Number, invalidPhoneNumber.CountryCode)

	assert.Equal("", formattedNumber)
	assert.NotNil(err)

	invalidPhoneNumber = phoneNumber{
		Number:      "255834643423323174",
		CountryCode: "GB",
	}

	formattedNumber, err = PhoneNumber(invalidPhoneNumber.Number, invalidPhoneNumber.CountryCode)

	assert.Equal("", formattedNumber)
	assert.NotNil(err)
}
