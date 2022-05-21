package app

//go:generate go run github.com/99designs/gqlgen
import (
	"github.com/go-playground/validator/v10"
	"gitlab.vecomentman.com/libs/logger"
	"gitlab.vecomentman.com/libs/sso"
	"gitlab.vecomentman.com/service/user/api"
	"gitlab.vecomentman.com/service/user/app/config"
)

type Resolver struct {
	ssoService     sso.Service
	userService    api.UserService
	companyService api.CompanyService
	validate       *validator.Validate
	config         *config.Config
	log            logger.Logger
}

func NewResolver(
	ssoService sso.Service,
	userService api.UserService,
	companyService api.CompanyService,
	validate *validator.Validate,
	config *config.Config,
	log logger.Logger,
) *Resolver {
	return &Resolver{
		ssoService:     ssoService,
		userService:    userService,
		companyService: companyService,
		validate:       validate,
		config:         config,
		log:            log,
	}
}
