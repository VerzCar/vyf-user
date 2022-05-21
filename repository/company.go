package repository

import (
	"context"
	"fmt"
	"gitlab.vecomentman.com/service/user/api/model"
	"gitlab.vecomentman.com/service/user/app/database"
	"gitlab.vecomentman.com/service/user/app/email"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CompanyById gets the country by id
func (s *storage) CompanyById(id int64) (*model.Company, error) {
	company := &model.Company{}
	err := s.db.Preload("Address.Country").
		Preload("Contact.PhoneNumberCountry").
		Preload("Contact.PhoneNumber2Country").
		Preload("Owner.Locale").
		Preload("Owner.Address").
		Preload("Owner.Address.Country").
		Preload("Owner.Contact").
		Preload("Owner.Contact.PhoneNumberCountry").
		Preload("Owner.Contact.PhoneNumber2Country").
		Preload(clause.Associations).Where(&model.Company{ID: id}).First(company).Error

	switch {
	case err != nil && !database.RecordNotFound(err):
		s.log.Errorf("error reading company by id %d: %s", id, err)
		return nil, err
	case database.RecordNotFound(err):
		s.log.Infof("company with id %d not found: %s", id, err)
		return nil, err
	}

	return company, nil
}

type SendCompanyVerificationEmailCallback func(*email.CompanyVerificationData) error
type RegisterNewCompanyCallback func(
	context.Context,
	string,
	int64,
) error

// CreateNewCompany based on given company model.
// Creates company in a transaction and call send email verification callback
// and register the company in the cache for verification.
// if one of the callbacks fails the transaction will do a rollback otherwise
// commit the transaction.
// Returns nil if no error occured otherwise an error
func (s *storage) CreateNewCompany(
	company *model.Company,
	sendCompanyVerificationEmail SendCompanyVerificationEmailCallback,
	emailData *email.CompanyVerificationData,
	registerNewCompany RegisterNewCompanyCallback,
	ctx context.Context,
	verificationKey string,
) error {
	err := s.db.Transaction(
		func(tx *gorm.DB) error {

			if err := tx.Create(company).Error; err != nil {
				return fmt.Errorf("error creating company: %s", err)
			}

			err := sendCompanyVerificationEmail(emailData)

			if err != nil {
				return fmt.Errorf("error sending verification email: %s", err)
			}

			err = registerNewCompany(ctx, verificationKey, company.ID)

			if err != nil {
				return err
			}

			return nil
		},
	)

	if err != nil {
		s.log.Error("error creating company: %s", err)
		return err
	}

	return nil
}

type CleanUpCompanyRegistrationCallback func(
	context.Context,
	string,
) error

// UpdateCompanyIsVerified update the company flag is verified in a transaction
// and call clean up of registration process for company callback
// if one of the callbacks fails the transaction will do a rollback otherwise
// commit the transaction.
// Returns nil if no error occured otherwise an error
func (s *storage) UpdateCompanyIsVerified(
	companyId int64,
	cleanUpCompanyRegistration CleanUpCompanyRegistrationCallback,
	ctx context.Context,
	verificationKey string,
) (*model.Company, error) {
	company, err := s.CompanyById(companyId)

	if err != nil {
		s.log.Error("getting company failed, error during reading from db: %s", err)
		return nil, err
	}

	err = s.db.Transaction(
		func(tx *gorm.DB) error {
			err = tx.Model(company).Update("IsVerified", true).Error

			if err != nil {
				s.log.Error("error updating company db entry: %s", err)
				return err
			}

			err = cleanUpCompanyRegistration(ctx, verificationKey)

			if err != nil {
				return err
			}

			return nil
		},
	)

	if err != nil {
		s.log.Error("error updating company: %s", err)
		return nil, err
	}

	return company, nil
}
