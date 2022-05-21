package repository

import (
	"gitlab.vecomentman.com/service/user/api/model"
	"gitlab.vecomentman.com/service/user/app/database"
)

// LocaleByLcidString gets the locale by lcid string. e.g.: "de-DE"
func (s *storage) LocaleByLcidString(lcid string) (*model.Locale, error) {
	locale := &model.Locale{}
	err := s.db.Where(&model.Locale{LcidString: lcid}).First(locale).Error

	switch {
	case err != nil && !database.RecordNotFound(err):
		s.log.Errorf("error reading locale by id %s: %s", lcid, err)
		return nil, err
	case database.RecordNotFound(err):
		s.log.Infof("locale with id %s not found: %s", lcid, err)
		return nil, err
	}

	return locale, nil
}
