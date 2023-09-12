package repository

import (
	"github.com/VerzCar/vyf-user/api/model"
	"github.com/VerzCar/vyf-user/validate"
)

// TransformContactRequest transform the input from source src to destination dest
func (s *storage) TransformContactRequest(src *model.ContactRequest, dest *model.Contact) error {

	dest.Email = src.Email

	country, err := s.CountryByAlpha2(src.PhoneNumberCountryAlphaCode)

	if err != nil {
		return err
	}

	e146PhoneNumber, err := validate.PhoneNumber(src.PhoneNumber, country.Alpha2)

	if err != nil {
		return err
	}

	dest.PhoneNumber = e146PhoneNumber
	dest.PhoneNumberCountry = country

	if src.PhoneNumber2CountryAlphaCode != nil && src.PhoneNumber2 != nil {

		country, err = s.CountryByAlpha2(*src.PhoneNumber2CountryAlphaCode)

		if err != nil {
			return err
		}

		e146PhoneNumber, err = validate.PhoneNumber(*src.PhoneNumber2, country.Alpha2)

		if err != nil {
			return err
		}

		dest.PhoneNumber2 = *src.PhoneNumber2
		dest.PhoneNumber2Country = country
	}

	if src.Web != nil {
		dest.Web = *src.Web
	}

	return nil
}
