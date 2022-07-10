package api

import (
	"context"
	"fmt"
	"gitlab.vecomentman.com/libs/logger"
	"gitlab.vecomentman.com/vote-your-face/service/user/api/model"
	"gitlab.vecomentman.com/vote-your-face/service/user/app/config"
	"gitlab.vecomentman.com/vote-your-face/service/user/app/database"
	emailSvc "gitlab.vecomentman.com/vote-your-face/service/user/app/email"
	routerContext "gitlab.vecomentman.com/vote-your-face/service/user/app/router/ctx"
)

type UserService interface {
	User(
		ctx context.Context,
	) (*model.User, error)
	UpdateUser(
		ctx context.Context,
		userInput *model.UserUpdateInput,
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

type userService struct {
	storage      UserRepository
	emailService emailSvc.Service
	config       *config.Config
	log          logger.Logger
}

func NewUserService(
	userRepo UserRepository,
	emailService emailSvc.Service,
	config *config.Config,
	log logger.Logger,
) UserService {
	return &userService{
		storage:      userRepo,
		emailService: emailService,
		config:       config,
		log:          log,
	}
}

func (u *userService) User(
	ctx context.Context,
) (*model.User, error) {
	authClaims, err := routerContext.ContextToAuthClaims(ctx)

	if err != nil {
		u.log.Errorf("error getting auth claims: %s", err)
		return nil, err
	}

	user, err := u.storage.UserByIdentityId(authClaims.Subject)

	switch {
	case err != nil && !database.RecordNotFound(err):
		u.log.Infof("could not query user for id: %s, error: %s", authClaims.Subject, err)
		return nil, err
	case database.RecordNotFound(err):
		newUser := &model.User{
			IdentityID: authClaims.Subject,
			Username:   "",
		}
		user, err := u.storage.CreateNewUser(newUser)

		if err != nil {
			u.log.Infof("could not create user for id: %s, error: %s", authClaims.Subject, err)
			return nil, err
		}

		return user, nil
	}

	return user, nil
}

func (u *userService) UpdateUser(
	ctx context.Context,
	userInput *model.UserUpdateInput,
) (*model.User, error) {
	user, err := u.User(ctx)

	if err != nil {
		return nil, err
	}

	if userInput.Username != nil {
		user.Username = *userInput.Username
	}

	if userInput.FirstName != nil {
		user.FirstName = *userInput.FirstName
	}

	if userInput.LastName != nil {
		user.LastName = *userInput.LastName
	}

	if userInput.Gender != nil {
		user.Gender = *userInput.Gender
	}

	if userInput.Locale != nil {
		locale, err := u.storage.LocaleByLcidString(*userInput.Locale)

		if err != nil {
			return nil, err
		}

		user.Locale = locale
	}

	if userInput.Profile != nil {
		profile := &model.Profile{}
		if user.Profile != nil {
			profile = user.Profile
		}
		transformProfileInput(userInput.Profile, profile)
		user.Profile = profile
	}

	if userInput.Address != nil {

		address := &model.Address{}
		if user.Address != nil {
			address = user.Address
		}
		err = u.storage.TransformAddressInput(userInput.Address, address)

		if err != nil {
			u.log.Errorf("error transforming address entry: %s", err)
			return nil, err
		}

		user.Address = address
	}

	if userInput.Contact != nil {

		contact := &model.Contact{}
		if user.Contact != nil {
			contact = user.Contact
		}
		err = u.storage.TransformContactInput(userInput.Contact, contact)

		if err != nil {
			u.log.Errorf("error transforming contact entry: %s", err)
			return nil, err
		}

		user.Contact = contact
	}

	user, err = u.storage.UpdateUser(user)

	if err != nil {
		return nil, fmt.Errorf("error updating user: %s", err)
	}

	return user, nil

}

func transformProfileInput(
	profileInput *model.ProfileInput,
	profile *model.Profile,
) {
	if profileInput.Bio != nil {
		profile.Bio = *profileInput.Bio
	}

	if profileInput.WhyVoteMe != nil {
		profile.WhyVoteMe = *profileInput.WhyVoteMe
	}

	if profileInput.ImageSrc != nil {
		profile.ImageSrc = *profileInput.ImageSrc
	}
}
