package validate

import (
	"fmt"
	"github.com/nyaruka/phonenumbers"
)

// PhoneNumber validates given phone numebr with its country code as region.
// If valid, it returns the formatted phone number in E164 format,
// otherwise an error.
func PhoneNumber(number string, countryAlphaCode string) (string, error) {

	parsedNumber, err := phonenumbers.Parse(number, countryAlphaCode)

	if err != nil {
		return "", fmt.Errorf("could not parse phone number %s with country code %s", number, countryAlphaCode)
	}

	if !phonenumbers.IsValidNumber(parsedNumber) {
		return "", fmt.Errorf("phone number %s not valid for region %s", number, countryAlphaCode)
	}

	formattedNumber := phonenumbers.Format(parsedNumber, phonenumbers.E164)

	return formattedNumber, nil
}
