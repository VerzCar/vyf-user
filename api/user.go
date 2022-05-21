package api

import (
	"context"
	"fmt"
	"gitlab.vecomentman.com/libs/logger"
	"gitlab.vecomentman.com/libs/sso"
	"gitlab.vecomentman.com/service/user/api/model"
	"gitlab.vecomentman.com/service/user/app/config"
	emailSvc "gitlab.vecomentman.com/service/user/app/email"
	routerContext "gitlab.vecomentman.com/service/user/app/router/ctx"
)

type UserService interface {
	CreateUser(
		ctx context.Context,
		userInput *model.UserCreateInput,
	) (*model.User, error)
	UpdateUser(
		ctx context.Context,
		user *model.UserUpdateInput,
	) (*model.User, error)
	User(
		ctx context.Context,
	) (*model.User, error)
}

type UserRepository interface {
	UserByIdentityId(id string) (*model.User, error)
	CreateNewUser(user *model.User) (*model.User, error)
	UpdateUser(user *model.User) (*model.User, error)
	LocaleByLcidString(lcid string) (*model.Locale, error)
	TransformAddressInput(src *model.AddressInput, dest *model.Address) error
	TransformContactInput(src *model.ContactInput, dest *model.Contact) error
}

type UserCache interface {
	StartResetUserPassword(
		ctx context.Context,
		passwordActivationKey string,
		userId string,
	) error
	UserInPasswordReset(
		ctx context.Context,
		resetPasswordKey string,
	) (string, error)
}

type userService struct {
	storage      UserRepository
	cache        UserCache
	ssoService   sso.Service
	emailService emailSvc.Service
	config       *config.Config
	log          logger.Logger
}

func NewUserService(
	userRepo UserRepository,
	cache UserCache,
	ssoService sso.Service,
	emailService emailSvc.Service,
	config *config.Config,
	log logger.Logger,
) UserService {
	return &userService{
		storage:      userRepo,
		cache:        cache,
		ssoService:   ssoService,
		emailService: emailService,
		config:       config,
		log:          log,
	}
}

func (u *userService) CreateUser(
	ctx context.Context,
	userInput *model.UserCreateInput,
) (*model.User, error) {

	identityId, err := u.ssoService.CreateUser(ctx, userInput.Email, userInput.Password)

	if err != nil {
		u.log.Error(err)
		return nil, err
	}

	user := &model.User{
		IdentityID: identityId,
	}

	user, err = u.storage.CreateNewUser(user)

	if err != nil {
		u.log.Errorf("could not create user, error: %s", err)

		if err := u.ssoService.DeleteUser(ctx, identityId); err != nil {
			u.log.Errorf("could not delete user, error: %s", err)
		}
		return nil, err
	}

	err = u.ssoService.SendVerificationEmail(ctx, identityId, "")

	if err != nil {
		u.log.Errorf("could not send verification email, error: %s", err)
	}

	return user, nil
}

func (u *userService) UpdateUser(
	ctx context.Context,
	user *model.UserUpdateInput,
) (*model.User, error) {
	currentUser, err := u.User(ctx)

	if err != nil {
		return nil, err
	}

	if user.FirstName != nil {
		currentUser.FirstName = *user.FirstName
	}

	if user.LastName != nil {
		currentUser.LastName = *user.LastName
	}

	if user.Gender != nil {
		currentUser.Gender = *user.Gender
	}

	if user.Locale != nil {
		locale, err := u.storage.LocaleByLcidString(*user.Locale)

		if err != nil {
			return nil, err
		}

		currentUser.Locale = locale
	}

	if user.AvatarURL != nil {
		currentUser.AvatarUrl = *user.AvatarURL
	}

	if user.Address != nil {

		address := &model.Address{}

		if currentUser.Address != nil {
			address = currentUser.Address
		}

		err := u.storage.TransformAddressInput(user.Address, address)

		if err != nil {
			u.log.Errorf("error transforming address entry: %s", err)
			return nil, err
		}

		currentUser.Address = address
	}

	if user.Contact != nil {

		contact := &model.Contact{}

		if currentUser.Contact != nil {
			contact = currentUser.Contact
		}

		ssoClaims, err := routerContext.ContextToSsoClaims(ctx)

		if err != nil {
			u.log.Errorf("error getting sso claims: %s", err)
		}

		// transform the ContactXsInput into the normal Contact model
		newContact := &model.ContactInput{
			Email:                        ssoClaims.Email,
			PhoneNumber:                  user.Contact.PhoneNumber,
			PhoneNumberCountryAlphaCode:  user.Contact.PhoneNumberCountryAlphaCode,
			PhoneNumber2:                 user.Contact.PhoneNumber2,
			PhoneNumber2CountryAlphaCode: user.Contact.PhoneNumber2CountryAlphaCode,
			Web:                          user.Contact.Web,
		}

		err = u.storage.TransformContactInput(newContact, contact)

		if err != nil {
			u.log.Errorf("error transforming contact entry: %s", err)
			return nil, err
		}

		currentUser.Contact = contact
	}

	currentUser, err = u.storage.UpdateUser(currentUser)

	if err != nil {
		return nil, fmt.Errorf("error updating user: %s", err)
	}

	return currentUser, nil
}

func (u *userService) User(
	ctx context.Context,
) (*model.User, error) {
	ssoClaims, err := routerContext.ContextToSsoClaims(ctx)

	if err != nil {
		u.log.Errorf("error getting sso claims: %s", err)
		return nil, err
	}

	user, err := u.storage.UserByIdentityId(ssoClaims.Subject)

	if err != nil {
		u.log.Infof("could not query user for id: %s, error: %s", ssoClaims.Subject, err)
		return nil, err
	}

	return user, nil
}
