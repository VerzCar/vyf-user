package app

import (
	awsx "github.com/VerzCar/vyf-lib-awsx"
	logger "github.com/VerzCar/vyf-lib-logger"
	"github.com/VerzCar/vyf-user/api"
	"github.com/VerzCar/vyf-user/app/config"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
)

type Server struct {
	router      *gin.Engine
	authService awsx.AuthService
	userService api.UserService
	validate    *validator.Validate
	config      *config.Config
	log         logger.Logger
}

func NewServer(
	router *gin.Engine,
	authService awsx.AuthService,
	userService api.UserService,
	validate *validator.Validate,
	config *config.Config,
	log logger.Logger,
) *Server {
	server := &Server{
		router:      router,
		authService: authService,
		userService: userService,
		validate:    validate,
		config:      config,
		log:         log,
	}

	server.routes()

	return server
}

func (s *Server) Run() error {
	err := s.router.Run(":8080")

	if err != nil {
		log.Printf("Server - there was an error calling Run on router: %v", err)
		return err
	}

	return nil
}
