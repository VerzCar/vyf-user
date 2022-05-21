package repository

import (
	"gitlab.vecomentman.com/service/user/api/model"
	"gitlab.vecomentman.com/service/user/app/database"
	"gorm.io/gorm/clause"
)

// UserById gets the user by id
func (s *storage) UserById(id int64) (*model.User, error) {
	user := &model.User{}
	err := s.db.Preload("Address.Country").
		Preload("Contact.PhoneNumberCountry").
		Preload("Contact.PhoneNumber2Country").
		Preload(clause.Associations).Where(&model.User{ID: id}).First(user).Error

	switch {
	case err != nil && !database.RecordNotFound(err):
		s.log.Errorf("error reading user by id %s: %s", id, err)
		return nil, err
	case database.RecordNotFound(err):
		s.log.Infof("user with id %s not found: %s", id, err)
		return nil, err
	}

	return user, nil
}

func (s *storage) UserByIdentityId(id string) (*model.User, error) {
	user := &model.User{}
	err := s.db.Preload("Address.Country").
		Preload("Contact.PhoneNumberCountry").
		Preload("Contact.PhoneNumber2Country").
		Preload(clause.Associations).Where(&model.User{IdentityID: id}).First(user).Error

	switch {
	case err != nil && !database.RecordNotFound(err):
		s.log.Errorf("error reading user by identity id %s: %s", id, err)
		return nil, err
	case database.RecordNotFound(err):
		s.log.Infof("user with identity id %s not found: %s", id, err)
		return nil, err
	}

	return user, nil
}

// CreateNewUser based on given user model
func (s *storage) CreateNewUser(user *model.User) (*model.User, error) {
	if err := s.db.Create(user).Error; err != nil {
		s.log.Infof("error creating user: %s", err)
		return nil, err
	}

	return user, nil
}

// UpdateUser update user based on given user model
func (s *storage) UpdateUser(user *model.User) (*model.User, error) {
	if err := s.db.Save(user).Error; err != nil {
		s.log.Errorf("error updating user: %s", err)
		return nil, err
	}

	return user, nil
}
