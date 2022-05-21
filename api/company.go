package api

import (
	"context"
	"fmt"
	"gitlab.vecomentman.com/libs/logger"
	"gitlab.vecomentman.com/service/user/api/model"
	"gitlab.vecomentman.com/service/user/app/config"
	"gitlab.vecomentman.com/service/user/app/email"
	"gitlab.vecomentman.com/service/user/repository"
	"gitlab.vecomentman.com/service/user/services"
	"gitlab.vecomentman.com/service/user/utils"
)

type CompanyService interface {
	CreateCompany(
		ctx context.Context,
		companyInput *model.CompanyCreateInput,
		currentUser *model.User,
	) (*model.Company, error)
	VerifyCompany(
		ctx context.Context,
		companyInput *model.CompanyVerifyInput,
	) (*model.Company, error)
}

type CompanyRepository interface {
	TransformAddressInput(src *model.AddressInput, dest *model.Address) error
	TransformContactInput(src *model.ContactInput, dest *model.Contact) error
	CreateNewCompany(
		company *model.Company,
		sendCompanyVerificationEmail repository.SendCompanyVerificationEmailCallback,
		emailData *email.CompanyVerificationData,
		registerNewCompany repository.RegisterNewCompanyCallback,
		ctx context.Context,
		verificationKey string,
	) error
	UpdateCompanyIsVerified(
		companyId int64,
		cleanUpCompanyRegistration repository.CleanUpCompanyRegistrationCallback,
		ctx context.Context,
		verificationKey string,
	) (*model.Company, error)
}

type CompanyCache interface {
	RegisterNewCompany(
		ctx context.Context,
		verificationKey string,
		companyId int64,
	) error
	CompanyInRegistration(
		ctx context.Context,
		verificationKey string,
	) (int64, error)
	CleanUpCompanyRegistration(
		ctx context.Context,
		verificationKey string,
	) error
}

type companyService struct {
	storage        CompanyRepository
	cache          CompanyCache
	emailService   email.Service
	paymentService services.PaymentService
	config         *config.Config
	log            logger.Logger
}

func NewCompanyService(
	companyRepo CompanyRepository,
	cache CompanyCache,
	emailService email.Service,
	paymentService services.PaymentService,
	config *config.Config,
	log logger.Logger,
) CompanyService {
	return &companyService{
		storage:        companyRepo,
		cache:          cache,
		emailService:   emailService,
		paymentService: paymentService,
		config:         config,
		log:            log,
	}
}

func (c *companyService) CreateCompany(
	ctx context.Context,
	companyInput *model.CompanyCreateInput,
	currentUser *model.User,
) (*model.Company, error) {
	// gets all the existing payments and check if payment exists
	// to create company
	payments, err := c.paymentService.Payments(ctx)

	if err != nil {
		c.log.Errorf("query payment service failed, error: %s", err)
		return nil, err
	}

	var paymentExist = false

	for _, payment := range payments {
		if payment.Id == companyInput.PaymentId {
			paymentExist = true
			break
		}
	}

	if paymentExist != true {
		c.log.Infof("payment missing")
		return nil, fmt.Errorf("payment missing")
	}

	// if payment exists create company
	address := &model.Address{}
	err = c.storage.TransformAddressInput(companyInput.Address, address)

	if err != nil {
		c.log.Errorf("error transforming address entry: %s", err)
		return nil, err
	}

	contact := &model.Contact{}

	err = c.storage.TransformContactInput(companyInput.Contact, contact)

	if err != nil {
		c.log.Errorf("error transforming contact entry: %s", err)
		return nil, err
	}

	companyUsers := []*model.User{currentUser}

	newCompany := &model.Company{
		Name:    companyInput.Name,
		Address: address,
		Contact: contact,
		Owner:   currentUser,
		Users:   companyUsers,
		TaxID:   companyInput.TaxID,
	}

	verificationKey := utils.UniqueKey()

	//verificationToken := token.NewJWToken(c.config.Security.Secrets.Key)
	//
	//verificationToken.ExpirationTimeDelta = c.config.Ttl.Token.Account.Verification
	//verificationToken.Claim.Subject = verificationKey
	//
	//err = verificationToken.Create()
	//
	//if err != nil {
	//	c.log.Errorf("error creating verification hash token: %s", err)
	//	return nil, err
	//}

	verificationEmailData := &email.CompanyVerificationData{
		VerificationToken:   "", //verificationToken.Signed,
		CustomerCompanyName: newCompany.Name,
		ToEmails:            []string{newCompany.Contact.Email},
	}

	err = c.storage.CreateNewCompany(
		newCompany,
		c.emailService.SendCompanyVerification,
		verificationEmailData,
		c.cache.RegisterNewCompany,
		ctx,
		verificationKey,
	)

	if err != nil {
		c.log.Error(err)
		return nil, err
	}

	return newCompany, nil
}

func (c *companyService) VerifyCompany(
	ctx context.Context,
	companyInput *model.CompanyVerifyInput,
) (*model.Company, error) {
	//verificationToken := token.NewJWToken(c.config.Security.Secrets.Key)
	//
	//verificationToken.Signed = companyInput.VerificationToken
	//
	//err := verificationToken.Verify()
	//
	//if err != nil {
	//	c.log.Infof("veryfing verification token failed: %s", err)
	//	return nil, err
	//}

	verifikationKey := "" //verificationToken.Claim.Subject

	companyId, err := c.cache.CompanyInRegistration(ctx, verifikationKey)

	if err != nil {
		return nil, err
	}

	company, err := c.storage.UpdateCompanyIsVerified(
		companyId,
		c.cache.CleanUpCompanyRegistration,
		ctx,
		verifikationKey,
	)

	if err != nil {
		return nil, err
	}

	return company, nil
}
