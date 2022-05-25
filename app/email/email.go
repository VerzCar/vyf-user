package email

import (
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"gitlab.vecomentman.com/libs/email"
	"gitlab.vecomentman.com/libs/logger"
	"gitlab.vecomentman.com/vote-your-face/service/user/app/config"
)

type Service interface {
	SendUserResetPasswordDone(data *UserResetPasswordDoneData) error
	SendUserResetPassword(data *UserResetPasswordData) error
	SendCompanyVerification(data *CompanyVerificationData) error
}

type service struct {
	config *config.Config
	log    logger.Logger
}

func NewService(
	config *config.Config,
	log logger.Logger,
) Service {
	return &service{
		config: config,
		log:    log,
	}
}

func (e *service) send(emailBlock email.Email) error {
	if e.config.Environment != config.EnvironmentDev {
		return emailBlock.Send()
	}
	return nil
}
