package email

import (
	"fmt"
	"gitlab.vecomentman.com/libs/email"
	"gitlab.vecomentman.com/vote-your-face/service/user/utils"
)

type CompanyVerificationData struct {
	VerificationToken   string
	CustomerCompanyName string
	ToEmails            []string
}

type companyVerificationEmailTemplateData struct {
	Header
	VerificationUrl     string
	CustomerCompanyName string
	Footer
}

func (e *service) SendCompanyVerification(data *CompanyVerificationData) error {
	layout := e.newLayoutData("Unternehmens-Verifizierung")

	verificationUrl := e.config.Hosts.Vec + "account/verification/company/" + data.VerificationToken

	templateData := companyVerificationEmailTemplateData{
		Header:              layout.Header,
		VerificationUrl:     verificationUrl,
		CustomerCompanyName: data.CustomerCompanyName,
		Footer:              layout.Footer,
	}

	emailTemplateHTMLPath := fmt.Sprintf("%semail-templates/dist/%s", utils.Base(), "verification.html")

	html, err := email.ParseHTMLTemplate(emailTemplateHTMLPath, templateData)

	if err != nil {
		return err
	}

	activationEmail := email.Email{
		Host: &email.Host{
			Name:     e.config.Smtp.Host,
			Port:     e.config.Smtp.Port,
			User:     e.config.Smtp.Account.User,
			Password: e.config.Smtp.Account.Password,
		},
		From: &email.From{
			Name:  e.config.Emails.From.Company.Name,
			Email: e.config.Emails.From.Company.Email,
		},
		To:      data.ToEmails,
		Subject: e.config.Emails.From.Company.Name,
		HTML:    html,
	}

	// send email here
	err = e.send(activationEmail)

	if err != nil {
		return err
	}

	return nil
}
