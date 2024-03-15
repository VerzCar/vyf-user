package repository

import (
	"github.com/VerzCar/vyf-user/api/model"
)

// TransformAddressRequest transform the input from source src to destination dest
func (s *storage) TransformAddressRequest(src *model.AddressRequest, dest *model.Address) error {
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
