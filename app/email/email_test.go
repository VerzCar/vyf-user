package email_test

import (
	"gitlab.vecomentman.com/vote-your-face/service/user/app/email"
	"reflect"
	"testing"
)

func TestService_SendCompanyVerification(t *testing.T) {
	emailSvc := email.NewService(config, log)

	data := &email.CompanyVerificationData{
		VerificationToken:   "923ujjdhijehw-34fkfkr3-34fk3ednjk",
		CustomerCompanyName: "Johns Pizza",
		ToEmails:            []string{"dev.anonymous@vecomentman.de"},
	}

	tests := []struct {
		name string
		data *email.CompanyVerificationData
		want error
	}{
		{
			name: "should send the email successfully",
			data: data,
			want: nil,
		},
	}

	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				err := emailSvc.SendCompanyVerification(test.data)

				if !reflect.DeepEqual(err, test.want) {
					t.Errorf("test: %v failed. \ngot: %v \nwanted: %v", test.name, err, test.want)
				}
			},
		)
	}
}

func TestService_SendUserResetPasswordDone(t *testing.T) {
	mockEmailSvc := email.NewService(config, log)

	data := &email.UserResetPasswordDoneData{
		ToEmails: []string{"dev.anonymous@vecomentman.de"},
	}

	tests := []struct {
		name string
		data *email.UserResetPasswordDoneData
		want error
	}{
		{
			name: "should send the email successfully",
			data: data,
			want: nil,
		},
	}

	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				err := mockEmailSvc.SendUserResetPasswordDone(test.data)

				if !reflect.DeepEqual(err, test.want) {
					t.Errorf("test: %v failed. \ngot: %v \nwanted: %v", test.name, err, test.want)
				}
			},
		)
	}
}

func TestService_SendUserResetPassword(t *testing.T) {
	mockEmailSvc := email.NewService(config, log)

	data := &email.UserResetPasswordData{
		ResetPasswordToken: "923ujjdhijehw-34fkfkr3-34fk3ednjk",
		ToEmails:           []string{"dev.anonymous@vecomentman.de"},
	}

	tests := []struct {
		name string
		data *email.UserResetPasswordData
		want error
	}{
		{
			name: "should send the email successfully",
			data: data,
			want: nil,
		},
	}

	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				err := mockEmailSvc.SendUserResetPassword(test.data)

				if !reflect.DeepEqual(err, test.want) {
					t.Errorf("test: %v failed. \ngot: %v \nwanted: %v", test.name, err, test.want)
				}
			},
		)
	}
}
