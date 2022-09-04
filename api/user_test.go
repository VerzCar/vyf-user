package api_test

import (
	"context"
	"database/sql"
	"fmt"
	"gitlab.vecomentman.com/vote-your-face/service/user/api"
	"gitlab.vecomentman.com/vote-your-face/service/user/api/model"
	"gitlab.vecomentman.com/vote-your-face/service/user/app/email"
	"gorm.io/gorm"
	"reflect"
	"testing"
	"time"
)

func TestUserService_User(t *testing.T) {
	mockRepo := &mockUserRepository{}
	mockEmailSvc := &mockEmailService{}
	mockService := api.NewUserService(mockRepo, mockEmailSvc, config, log)

	ctxMock := putUserIntoContext(mockUser.Elon)
	ctxNewUserMock := putUserIntoContext(mockUser.NewUser)

	anonymousUser := "newUser"

	tests := []struct {
		name       string
		ctx        context.Context
		identityId *string
		want       error
	}{
		{
			name:       "should query user successfully",
			ctx:        ctxMock,
			identityId: nil,
			want:       nil,
		},
		{
			name:       "should query user first time successfully",
			ctx:        ctxNewUserMock,
			identityId: nil,
			want:       nil,
		},
		{
			name:       "should fail because given user id does not exist in db",
			ctx:        ctxNewUserMock,
			identityId: &anonymousUser,
			want:       fmt.Errorf("record not found"),
		},
		{
			name:       "should fail because user is not authenticated",
			ctx:        emptyUserContext(),
			identityId: nil,
			want:       fmt.Errorf("could not retrieve auth claims"),
		},
	}

	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				_, err := mockService.User(test.ctx, test.identityId)

				if !reflect.DeepEqual(err, test.want) {
					t.Errorf("test: %v failed. \ngot: %v \nwanted: %v", test.name, err, test.want)
				}
			},
		)
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	mockRepo := &mockUserRepository{}
	mockEmailSvc := &mockEmailService{}
	mockService := api.NewUserService(mockRepo, mockEmailSvc, config, log)

	ctxMock := putUserIntoContext(mockUser.Elon)
	//ctxNewUserMock := putUserIntoContext(mockUser.NewUser)

	firstName := "firstName"
	lastName := "lastName"
	username := "username"
	gender := model.GenderWomen
	locale := "de-DE"
	profile := model.ProfileInput{
		Bio:       stringP("bio"),
		WhyVoteMe: stringP("whyVoteMe"),
		ImageSrc:  stringP("imageSrc"),
	}
	input := model.UserUpdateInput{
		FirstName: &firstName,
		LastName:  &lastName,
		Username:  &username,
		Gender:    &gender,
		Locale:    &locale,
		Address:   &model.AddressInput{},
		Contact:   &model.ContactInput{},
		Profile:   &profile,
	}

	tests := []struct {
		name  string
		ctx   context.Context
		input *model.UserUpdateInput
		want  error
	}{
		{
			name:  "should update user successfully",
			ctx:   ctxMock,
			input: &input,
			want:  nil,
		},
		//{
		//	name:       "should query user first time successfully",
		//	ctx:        ctxNewUserMock,
		//	identityId: nil,
		//	want:       nil,
		//},
		//{
		//	name:       "should fail because given user id does not exist in db",
		//	ctx:        ctxNewUserMock,
		//	identityId: &anonymousUser,
		//	want:       fmt.Errorf("record not found"),
		//},
		{
			name:  "should fail because user is not authenticated",
			ctx:   emptyUserContext(),
			input: &input,
			want:  fmt.Errorf("could not retrieve auth claims"),
		},
	}

	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				_, err := mockService.UpdateUser(test.ctx, test.input)

				if !reflect.DeepEqual(err, test.want) {
					t.Errorf("test: %v failed. \ngot: %v \nwanted: %v", test.name, err, test.want)
				}
			},
		)
	}
}

type mockUserRepository struct{}

type mockEmailService struct{}

func (m mockUserRepository) UserByIdentityId(id string) (*model.User, error) {
	user := &model.User{
		ID:         1,
		IdentityID: id,
		Username:   "elon",
		FirstName:  "elon",
		LastName:   "musk",
		Gender:     model.GenderMan,
		LocaleID:   sql.NullInt64{},
		Locale:     nil,
		AddressID:  sql.NullInt64{},
		Address:    nil,
		ContactID:  sql.NullInt64{},
		Contact:    nil,
		ProfileID:  sql.NullInt64{},
		Profile:    &model.Profile{},
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	switch id {
	case "newUser":
		return nil, gorm.ErrRecordNotFound
	}

	return user, nil
}

func (m mockUserRepository) CreateNewUser(user *model.User) (*model.User, error) {
	return &model.User{}, nil
}

func (m mockUserRepository) UpdateUser(user *model.User) (*model.User, error) {
	return &model.User{}, nil
}

func (m mockUserRepository) LocaleByLcidString(lcid string) (*model.Locale, error) {
	return &model.Locale{}, nil
}

func (m mockUserRepository) TransformAddressInput(src *model.AddressInput, dest *model.Address) error {
	return nil
}

func (m mockUserRepository) TransformContactInput(src *model.ContactInput, dest *model.Contact) error {
	return nil
}

func (m mockEmailService) SendUserResetPasswordDone(data *email.UserResetPasswordDoneData) error {
	//TODO implement me
	panic("implement me")
}

func (m mockEmailService) SendUserResetPassword(data *email.UserResetPasswordData) error {
	//TODO implement me
	panic("implement me")
}

func stringP(val string) *string {
	return &val
}
