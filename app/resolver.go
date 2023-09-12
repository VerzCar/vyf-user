package app

//go:generate go run github.com/99designs/gqlgen
import (
	"github.com/VerzCar/vyf-lib-awsx"
	"github.com/VerzCar/vyf-lib-logger"
	"github.com/VerzCar/vyf-user/api"
	"github.com/VerzCar/vyf-user/app/config"
	"github.com/go-playground/validator/v10"
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
