package app

import (
	"fmt"
	awsx "github.com/VerzCar/vyf-lib-awsx"
	logger "github.com/VerzCar/vyf-lib-logger"
	"github.com/VerzCar/vyf-user/api"
	"github.com/VerzCar/vyf-user/app/config"
	"github.com/VerzCar/vyf-user/app/sanitizer"
	"github.com/gin-gonic/gin"
	"log"
)

type Server struct {
	router            *gin.Engine
	authService       awsx.AuthService
	userService       api.UserService
	userUploadService api.UserUploadService
	validate          sanitizer.Validator
	config            *config.Config
	log               logger.Logger
}

func NewServer(
	router *gin.Engine,
	authService awsx.AuthService,
	userService api.UserService,
	userUploadService api.UserUploadService,
	validate sanitizer.Validator,
	config *config.Config,
	log logger.Logger,
) *Server {
	server := &Server{
		router:            router,
		authService:       authService,
		userService:       userService,
		userUploadService: userUploadService,
		validate:          validate,
		config:            config,
		log:               log,
	}

	server.routes()

	return server
}

func (s *Server) Run() error {
	port := fmt.Sprintf(":%s", s.config.Port)
	err := s.router.Run(port)

	if err != nil {
		log.Printf("Server - there was an error calling Run on router: %v", err)
		return err
	}

	return nil
}
