package api

import (
	"context"
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
	emailService emailSvc.Service
	config       *config.Config
	log          logger.Logger
}

func NewUserService(
	userRepo UserRepository,
	cache UserCache,
	emailService emailSvc.Service,
	config *config.Config,
	log logger.Logger,
) UserService {
	return &userService{
		storage:      userRepo,
		cache:        cache,
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
