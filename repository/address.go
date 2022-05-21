package repository

import (
	"gitlab.vecomentman.com/service/user/api/model"
)

// TransformAddressInput transform the input from source src to destination dest
func (s *storage) TransformAddressInput(src *model.AddressInput, dest *model.Address) error {
	dest.Address = src.Address
	dest.City = src.City
	dest.PostalCode = src.PostalCode

	country, err := s.CountryByAlpha2(src.CountryAlphaCode)

	if err != nil {
		return err
	}

	dest.Country = country

	return nil
}
