package email

import (
	"fmt"
	"gitlab.vecomentman.com/libs/email"
	"gitlab.vecomentman.com/vote-your-face/service/user/utils"
)

type UserResetPasswordDoneData struct {
	ToEmails []string
}

type userResetPasswordDoneEmailTemplateData struct {
	Header
	Footer
}

func (e *service) SendUserResetPasswordDone(data *UserResetPasswordDoneData) error {
	layout := e.newLayoutData("Passwort geändert")

	templateData := userResetPasswordDoneEmailTemplateData{
		Header: layout.Header,
		Footer: layout.Footer,
	}

	emailTemplateHTMLPath := fmt.Sprintf("%semail-templates/dist/%s", utils.Base(), "password_reset_done.html")

	html, err := email.ParseHTMLTemplate(emailTemplateHTMLPath, templateData)

	if err != nil {
		return err
	}

	resetPasswordEmail := email.Email{
		Host: &email.Host{
			Name:     e.config.Smtp.Host,
			Port:     e.config.Smtp.Port,
			User:     e.config.Smtp.NoReply.User,
			Password: e.config.Smtp.NoReply.Password,
		},
		From: &email.From{
			Name:  e.config.Emails.From.NoReply.Name,
			Email: e.config.Emails.From.NoReply.Email,
		},
		To:      data.ToEmails,
		Subject: e.config.Emails.From.NoReply.Name,
		HTML:    html,
	}

	err = e.send(resetPasswordEmail)

	if err != nil {
		return err
	}

	return nil
}

type UserResetPasswordData struct {
	ResetPasswordToken string
	ToEmails           []string
}

type userResetPasswordEmailTemplateData struct {
	Header
	ResetPasswordUrl string
	Footer
}

func (e *service) SendUserResetPassword(data *UserResetPasswordData) error {
	layout := e.newLayoutData("Passwort zurücksetzen")

	activationUrl := e.config.Hosts.Vec + "account/reset-password/" + data.ResetPasswordToken

	templateData := userResetPasswordEmailTemplateData{
		Header:           layout.Header,
		ResetPasswordUrl: activationUrl,
		Footer:           layout.Footer,
	}

	emailTemplateHTMLPath := fmt.Sprintf("%semail-templates/dist/%s", utils.Base(), "reset_password.html")

	html, err := email.ParseHTMLTemplate(emailTemplateHTMLPath, templateData)

	if err != nil {
		return err
	}

	resetPasswordEmail := email.Email{
		Host: &email.Host{
			Name:     e.config.Smtp.Host,
			Port:     e.config.Smtp.Port,
			User:     e.config.Smtp.NoReply.User,
			Password: e.config.Smtp.NoReply.Password,
		},
		From: &email.From{
			Name:  e.config.Emails.From.NoReply.Name,
			Email: e.config.Emails.From.NoReply.Email,
		},
		To:      data.ToEmails,
		Subject: e.config.Emails.From.NoReply.Name,
		HTML:    html,
	}

	err = e.send(resetPasswordEmail)

	if err != nil {
		return err
	}

	return nil
}
