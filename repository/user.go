package repository

import (
	"fmt"
	"github.com/VerzCar/vyf-user/api/model"
	"github.com/VerzCar/vyf-user/app/database"
	"gorm.io/gorm"
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

func (s *storage) Users(identityID string) ([]*model.UserPaginated, error) {
	var users []*model.UserPaginated

	err := s.db.Model(&model.User{}).
		Select("users.username, users.identity_id, profiles.image_src as profile_image_src").
		Joins("left join profiles on profiles.id = users.profile_id").
		Not(&model.User{IdentityID: identityID}).
		Limit(100).
		Order("username ::bytea").
		Find(&users).
		Error

	switch {
	case err != nil && !database.RecordNotFound(err):
		s.log.Errorf("error reading users: %s", err)
		return nil, err
	case database.RecordNotFound(err):
		s.log.Infof("users not found: %s", err)
		return nil, err
	}

	return users, nil
}

func (s *storage) UsersFiltered(
	username string,
) ([]*model.UserPaginated, error) {
	var users []*model.UserPaginated

	err := s.db.Model(&model.User{}).
		Select("users.username, users.identity_id, profiles.image_src as profile_image_src").
		Joins("left join profiles on profiles.id = users.profile_id").
		Limit(100).
		Order("username ::bytea").
		Where("username ILIKE ?", fmt.Sprintf("%%%s%%", username)).
		Find(&users).
		Error

	switch {
	case err != nil && !database.RecordNotFound(err):
		s.log.Errorf("error reading users: %s", err)
		return nil, err
	case database.RecordNotFound(err):
		s.log.Infof("users not found: %s", err)
		return nil, err
	}

	return users, nil
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
	if err := s.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(user).Error; err != nil {
		s.log.Errorf("error updating user: %s", err)
		return nil, err
	}

	return user, nil
}
