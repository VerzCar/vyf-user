package api

import (
	"context"
	"fmt"
	"gitlab.vecomentman.com/libs/logger"
	"gitlab.vecomentman.com/vote-your-face/service/user/api/model"
	"gitlab.vecomentman.com/vote-your-face/service/user/app/config"
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

	fmt.Println(authClaims.Subject)

	//user, err := u.storage.UserByIdentityId(ssoClaims.Subject)
	//
	//if err != nil {
	//	u.log.Infof("could not query user for id: %s, error: %s", ssoClaims.Subject, err)
	//	return nil, err
	//}

	user := &model.User{IdentityID: "yeah", FirstName: "Carlo"}

	return user, nil
}
