package app

//go:generate go run github.com/99designs/gqlgen
import (
	"github.com/go-playground/validator/v10"
	"gitlab.vecomentman.com/libs/awsx"
	"gitlab.vecomentman.com/libs/logger"
	"gitlab.vecomentman.com/vote-your-face/service/user/api"
	"gitlab.vecomentman.com/vote-your-face/service/user/app/config"
)

type Resolver struct {
	authService awsx.AuthService
	userService api.UserService
	validate    *validator.Validate
	config      *config.Config
	log         logger.Logger
}

func NewResolver(
	authService awsx.AuthService,
	userService api.UserService,
	validate *validator.Validate,
	config *config.Config,
	log logger.Logger,
) *Resolver {
	return &Resolver{
		authService: authService,
		userService: userService,
		validate:    validate,
		config:      config,
		log:         log,
	}
}
