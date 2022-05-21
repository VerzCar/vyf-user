package repository

import (
	"gitlab.vecomentman.com/service/user/api/model"
	"gitlab.vecomentman.com/service/user/app/database"
)

// CountryById gets the country by id
func (s *storage) CountryById(id int64) (*model.Country, error) {
	country := &model.Country{}
	err := s.db.Where(&model.Country{ID: id}).First(country).Error

	switch {
	case err != nil && !database.RecordNotFound(err):
		s.log.Errorf("error reading country by id %d: %s", id, err)
		return nil, err
	case database.RecordNotFound(err):
		s.log.Infof("country with id %d not found: %s", id, err)
		return nil, err
	}

	return country, nil
}

// CountryByAlpha2 gets the country by alpha2 code
func (s *storage) CountryByAlpha2(alpha2 string) (*model.Country, error) {
	country := &model.Country{}
	err := s.db.Where(&model.Country{Alpha2: alpha2}).First(country).Error

	switch {
	case err != nil && !database.RecordNotFound(err):
		s.log.Errorf("error reading country by id %s: %s", alpha2, err)
		return nil, err
	case database.RecordNotFound(err):
		s.log.Infof("country with id %s not found: %s", alpha2, err)
		return nil, err
	}

	return country, nil
}
