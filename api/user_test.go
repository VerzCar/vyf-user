package api_test

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"gitlab.vecomentman.com/libs/sso"
	"gitlab.vecomentman.com/vote-your-face/service/user/api"
	"gitlab.vecomentman.com/vote-your-face/service/user/api/model"
	"gitlab.vecomentman.com/vote-your-face/service/user/app/email"
	"gitlab.vecomentman.com/vote-your-face/service/user/test/factory"
	"gitlab.vecomentman.com/vote-your-face/service/user/validate"
	"gorm.io/gorm"
	"reflect"
	"testing"
)

type mockUserRepo struct{}

type mockUserCache struct{}

type mockEmailService struct{}

type mockSsoService struct{}

type ssoStorages struct {
	SsoElon sso.Claims
}

type userStorages struct {
	Elon model.User
}

var ssoStorage = ssoStorages{
	SsoElon: factory.NewSsoUsers().Elon,
}

var userStorage = userStorages{
	Elon: factory.NewCreationUsers().Elon,
}
var hasCreatedUserElonInSso = false

func (s mockSsoService) DecodeAccessToken(ctx context.Context, accessToken string, realm string) (
	*jwt.Token,
	*sso.Claims,
	error,
) {
	panic("implement me")
}

func (s mockSsoService) DefaultRealm() string {
	return config.Sso.Realm.Name
}

func (s mockSsoService) CreateUser(ctx context.Context, email string, password string) (string, error) {

	if ssoStorage.SsoElon.Email == email && hasCreatedUserElonInSso {
		return "", fmt.Errorf("user with username already exists")
	}

	hasCreatedUserElonInSso = true
	return ssoStorage.SsoElon.Subject, nil
}

func (s mockSsoService) DeleteUser(ctx context.Context, userId string) error {
	if ssoStorage.SsoElon.Subject == userId {
		ssoStorage.SsoElon = sso.Claims{}
		hasCreatedUserElonInSso = false
		return nil
	}

	return fmt.Errorf("user with user id does not exist")
}

func (s mockSsoService) SendVerificationEmail(ctx context.Context, userId string) error {
	return nil
}

func (e mockEmailService) SendUserResetPasswordDone(data *email.UserResetPasswordDoneData) error {
	return nil
}

func (e mockEmailService) SendUserResetPassword(data *email.UserResetPasswordData) error {
	return nil
}

func (e mockEmailService) SendCompanyVerification(data *email.CompanyVerificationData) error {
	return nil
}

func (m mockUserRepo) UserById(id int64) (*model.User, error) {
	panic("implement me")
}

func (m mockUserRepo) UserByIdentityId(id string) (*model.User, error) {
	switch id {
	case dbUser.Martin.IdentityID:
		return &dbUser.Martin, nil
	default:
		return nil, gorm.ErrRecordNotFound
	}
}

func (m mockUserRepo) CreateNewUser(user *model.User) (*model.User, error) {
	userStorage.Elon.ID = 555
	userStorage.Elon.IdentityID = ssoStorage.SsoElon.Subject
	return &userStorage.Elon, nil
}

func (m mockUserRepo) UpdateUser(user *model.User) (*model.User, error) {
	return user, nil
}

func (m mockUserRepo) TransformAddressInput(src *model.AddressInput, dest *model.Address) error {
	dest.Address = src.Address
	dest.City = src.City
	dest.PostalCode = src.PostalCode

	var country *model.Country
	var err error

	switch src.CountryAlphaCode {
	case factory.CountryGermany.Alpha2:
		country = factory.CountryGermany
	default:
		err = gorm.ErrRecordNotFound
	}

	if err != nil {
		return err
	}

	dest.Country = country

	return nil
}

func (m mockUserRepo) TransformContactInput(src *model.ContactInput, dest *model.Contact) error {
	dest.Email = src.Email

	var country *model.Country
	var err error

	switch src.PhoneNumberCountryAlphaCode {
	case factory.CountryGermany.Alpha2:
		country = factory.CountryGermany
	default:
		err = gorm.ErrRecordNotFound
	}

	if err != nil {
		return err
	}

	e146PhoneNumber, err := validate.PhoneNumber(src.PhoneNumber, country.Alpha2)

	if err != nil {
		return err
	}

	dest.PhoneNumber = e146PhoneNumber
	dest.PhoneNumberCountry = country

	if src.PhoneNumber2CountryAlphaCode != nil && src.PhoneNumber2 != nil {

		switch *src.PhoneNumber2CountryAlphaCode {
		case factory.CountryGermany.Alpha2:
			country = factory.CountryGermany
		default:
			err = gorm.ErrRecordNotFound
		}

		if err != nil {
			return err
		}

		e146PhoneNumber, err = validate.PhoneNumber(*src.PhoneNumber2, country.Alpha2)

		if err != nil {
			return err
		}

		dest.PhoneNumber2 = *src.PhoneNumber2
		dest.PhoneNumber2Country = country
	}

	if src.Web != nil {
		dest.Web = *src.Web
	}

	return nil
}

func (m mockUserRepo) LocaleByLcidString(lcid string) (*model.Locale, error) {
	switch lcid {
	case factory.LocaleGermany.LcidString:
		return factory.LocaleGermany, nil
	default:
		return nil, gorm.ErrRecordNotFound
	}
}

func (m mockUserCache) StartResetUserPassword(
	ctx context.Context,
	passwordActivationKey string,
	userId string,
) error {
	panic("implement me")
}
func (m mockUserCache) UserInPasswordReset(
	ctx context.Context,
	resetPasswordKey string,
) (string, error) {
	panic("implement me")
}

func TestUserService_CreateUser(t *testing.T) {
	mockRepo := &mockUserRepo{}
	mockCache := &mockUserCache{}
	mockEmailSvc := &mockEmailService{}
	mockSsoSvc := &mockSsoService{}
	mockUserService := api.NewUserService(mockRepo, mockCache, mockSsoSvc, mockEmailSvc, config, log)

	user := factory.NewInputUsers()

	tests := []struct {
		name string
		ctx  context.Context
		user *model.UserCreateInput
		want error
	}{
		{
			name: "should create a new user successfully",
			ctx:  ctx,
			user: &user.Elon,
			want: nil,
		},
		{
			name: "should fail because user with same username already exists",
			ctx:  ctx,
			user: &user.Elon,
			want: fmt.Errorf("user with username already exists"),
		},
	}

	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				_, err := mockUserService.CreateUser(test.ctx, test.user)

				if !reflect.DeepEqual(err, test.want) {
					t.Errorf("test: %v failed. \ngot: %v \nwanted: %v", test.name, err, test.want)
				}
			},
		)
	}
}

func TestUserService_User(t *testing.T) {
	mockRepo := &mockUserRepo{}
	mockCache := &mockUserCache{}
	mockEmailSvc := &mockEmailService{}
	mockSsoSvc := &mockSsoService{}
	mockUserService := api.NewUserService(mockRepo, mockCache, mockSsoSvc, mockEmailSvc, config, log)

	ssoMartin := factory.NewSsoUsers().Martin
	ssoElon := factory.NewSsoUsers().Elon

	tests := []struct {
		name         string
		ctx          context.Context
		expectedUser *model.User
		want         error
	}{
		{
			name:         "should return the logged in user",
			ctx:          ssoUserIntoContext(&ssoMartin),
			expectedUser: &dbUser.Martin,
			want:         nil,
		},
		{
			name:         "should fail because logged in user does not exists on database",
			ctx:          ssoUserIntoContext(&ssoElon),
			expectedUser: nil,
			want:         fmt.Errorf("record not found"),
		},
		{
			name:         "should fail because sso claims are not set into context",
			ctx:          ssoUserIntoContextFail(nil),
			expectedUser: nil,
			want:         fmt.Errorf("could not retrieve sso claims"),
		},
	}

	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				user, err := mockUserService.User(test.ctx)

				if !reflect.DeepEqual(err, test.want) {
					t.Errorf("test: %v failed. \ngot: %v \nwanted: %v", test.name, err, test.want)
				}

				if err == nil && !reflect.DeepEqual(user, test.expectedUser) {
					t.Errorf("test: %v failed. \nuser: %v \nwanted: %v", test.name, user, test.expectedUser)
				}
			},
		)
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	mockRepo := &mockUserRepo{}
	mockCache := &mockUserCache{}
	mockEmailSvc := &mockEmailService{}
	mockSsoSvc := &mockSsoService{}
	mockUserService := api.NewUserService(mockRepo, mockCache, mockSsoSvc, mockEmailSvc, config, log)

	ssoMartin := factory.NewSsoUsers().Martin

	userInput := factory.NewUpdateInputUsers().Martin

	tests := []struct {
		name        string
		ctx         context.Context
		user        *model.UserUpdateInput
		updateValue string
		want        error
	}{
		{
			name:        "should update the user successfully",
			ctx:         ssoUserIntoContext(&ssoMartin),
			user:        &userInput,
			updateValue: *userInput.FirstName,
			want:        nil,
		},
		{
			name:        "should fail because sso claims are not set into context",
			ctx:         ssoUserIntoContextFail(nil),
			user:        &userInput,
			updateValue: *userInput.FirstName,
			want:        fmt.Errorf("could not retrieve sso claims"),
		},
	}

	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				updatedUser, err := mockUserService.UpdateUser(test.ctx, test.user)

				if !reflect.DeepEqual(err, test.want) {
					t.Errorf("test: %v failed. \ngot: %v \nwanted: %v", test.name, err, test.want)
					return
				}

				if err == nil && !reflect.DeepEqual(updatedUser.FirstName, test.updateValue) {
					t.Errorf(
						"test: %v failed. \nupdated user: %v \nnot equal desired: %v",
						test.name,
						updatedUser.FirstName,
						test.updateValue,
					)
				}
			},
		)
	}
}
