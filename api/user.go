package api

import (
	"bytes"
	"context"
	"fmt"
	"github.com/VerzCar/vyf-lib-logger"
	"github.com/VerzCar/vyf-user/api/model"
	"github.com/VerzCar/vyf-user/app/config"
	"github.com/VerzCar/vyf-user/app/database"
	routerContext "github.com/VerzCar/vyf-user/app/router/ctx"
	"net/http"
	"strings"
	"time"
)

type UserService interface {
	User(
		ctx context.Context,
		identityId *string,
	) (*model.User, error)
	Users(
		ctx context.Context,
	) ([]*model.UserPaginated, error)
	UsersFiltered(
		ctx context.Context,
		username *string,
	) ([]*model.UserPaginated, error)
	UpdateUser(
		ctx context.Context,
		userInput *model.UserUpdateRequest,
	) (*model.User, error)
}

type UserRepository interface {
	UserByIdentityId(id string) (*model.User, error)
	Users(identityID string) ([]*model.UserPaginated, error)
	UsersFiltered(
		callerIdentityID string,
		username string,
	) ([]*model.UserPaginated, error)
	CreateNewUser(user *model.User) (*model.User, error)
	UpdateUser(user *model.User) (*model.User, error)
	LocaleByLcidString(lcid string) (*model.Locale, error)
	TransformAddressRequest(src *model.AddressRequest, dest *model.Address) error
	TransformContactRequest(src *model.ContactRequest, dest *model.Contact) error
}

type userService struct {
	storage UserRepository
	config  *config.Config
	log     logger.Logger
}

func NewUserService(
	userRepo UserRepository,
	config *config.Config,
	log logger.Logger,
) UserService {
	return &userService{
		storage: userRepo,
		config:  config,
		log:     log,
	}
}

// User from the user repository.
// Gets the user if it exists. If it does not exist the user will be created the first time.
// If the identityId is provided the queried user will be returned if exists.
func (u *userService) User(
	ctx context.Context,
	identityId *string,
) (*model.User, error) {
	authClaims, err := routerContext.ContextToAuthClaims(ctx)

	if err != nil {
		u.log.Errorf("error getting auth claims: %s", err)
		return nil, err
	}

	var queryIdentityId string
	isQueryingItself := identityId == nil

	if isQueryingItself {
		queryIdentityId = authClaims.Subject
	} else {
		queryIdentityId = *identityId
	}

	user, err := u.storage.UserByIdentityId(queryIdentityId)

	switch {
	case err != nil && !database.RecordNotFound(err):
		u.log.Infof("could not query user for id: %s, error: %s", queryIdentityId, err)
		return nil, err
	case database.RecordNotFound(err) && isQueryingItself:
		newUser := &model.User{
			IdentityID: queryIdentityId,
			Username:   authClaims.PrivateClaims.Username,
			Profile:    &model.Profile{},
		}
		user, err := u.storage.CreateNewUser(newUser)

		if err != nil {
			u.log.Infof("could not create user for id: %s, error: %s", queryIdentityId, err)
			return nil, err
		}

		err = addUserToGlobalCircle(ctx, u.config.Host.Service.VoteCircle)

		if err != nil {
			u.log.Errorf("could not add user to global circle for user id: %s, error: %s", queryIdentityId, err)
		}

		return user, nil
	case database.RecordNotFound(err):
		u.log.Infof("could not query user for id: %s, error: %s", queryIdentityId, err)
		return nil, err
	}

	return user, nil
}

func (u *userService) UpdateUser(
	ctx context.Context,
	userInput *model.UserUpdateRequest,
) (*model.User, error) {
	user, err := u.User(ctx, nil)

	if err != nil {
		return nil, err
	}

	if userInput.Username != nil {
		user.Username = strings.TrimSpace(*userInput.Username)
	}

	if userInput.FirstName != nil {
		user.FirstName = strings.TrimSpace(*userInput.FirstName)
	}

	if userInput.LastName != nil {
		user.LastName = strings.TrimSpace(*userInput.LastName)
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
		err = u.storage.TransformAddressRequest(userInput.Address, address)

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
		err = u.storage.TransformContactRequest(userInput.Contact, contact)

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

func (u *userService) Users(
	ctx context.Context,
) ([]*model.UserPaginated, error) {
	authClaims, err := routerContext.ContextToAuthClaims(ctx)

	if err != nil {
		u.log.Errorf("error getting auth claims: %s", err)
		return nil, err
	}

	users, err := u.storage.Users(authClaims.Subject)

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *userService) UsersFiltered(
	ctx context.Context,
	username *string,
) ([]*model.UserPaginated, error) {
	authClaims, err := routerContext.ContextToAuthClaims(ctx)

	if err != nil {
		u.log.Errorf("error getting auth claims: %s", err)
		return nil, err
	}

	users, err := u.storage.UsersFiltered(authClaims.Subject, *username)

	if err != nil {
		return nil, err
	}

	return users, nil
}

func transformProfileInput(
	profileInput *model.ProfileRequest,
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

// Add the user to the global circle if it is a new user.
// Makes a http PUT request to vote-circle service.
func addUserToGlobalCircle(
	ctx context.Context,
	serviceUrl string,
) error {
	accessToken, err := routerContext.ContextToBearerToken(ctx)

	if err != nil {
		return err
	}

	jsonBody := []byte(``)
	bodyReader := bytes.NewReader(jsonBody)

	requestURL := fmt.Sprintf("%s/circle/to-global", serviceUrl)
	req, err := http.NewRequest(http.MethodPut, requestURL, bodyReader)

	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	_, err = client.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		return err
	}

	return nil
}
