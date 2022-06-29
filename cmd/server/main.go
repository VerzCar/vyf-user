package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"gitlab.vecomentman.com/libs/awsx"
	"gitlab.vecomentman.com/libs/logger"
	"gitlab.vecomentman.com/vote-your-face/service/user/api"
	"gitlab.vecomentman.com/vote-your-face/service/user/app"
	"gitlab.vecomentman.com/vote-your-face/service/user/app/config"
	"gitlab.vecomentman.com/vote-your-face/service/user/app/database"
	"gitlab.vecomentman.com/vote-your-face/service/user/app/email"
	"gitlab.vecomentman.com/vote-your-face/service/user/app/router"
	"gitlab.vecomentman.com/vote-your-face/service/user/repository"
	"gitlab.vecomentman.com/vote-your-face/service/user/utils"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "startup error: %s\\n", err)
		os.Exit(1)
	}
}

var validate *validator.Validate

func run() error {
	configPath := utils.FromBase("app/config/")

	envConfig := config.NewConfig(configPath)
	log := logger.NewLogger(configPath)

	log.Infof("Startup service...")
	log.Infof("Configuration loaded.")

	db := database.Connect(log, envConfig)

	storage := repository.NewStorage(db, envConfig, log)

	sqlDb, _ := db.DB()
	err := storage.RunMigrationsUp(sqlDb)

	if err != nil {
		return err
	}

	// initialize auth service
	authService, err := awsx.NewAuthService(
		awsx.AppClientId(envConfig.Aws.Auth.ClientId),
		awsx.ClientSecret(envConfig.Aws.Auth.ClientSecret),
		awsx.AwsDefaultRegion(envConfig.Aws.Auth.AwsDefaultRegion),
		awsx.UserPoolId(envConfig.Aws.Auth.UserPoolId),
	)

	if err != nil {
		return err
	}

	// initialize third party services
	emailService := email.NewService(envConfig, log)

	if err != nil {
		return err
	}

	// initialize api services
	userService := api.NewUserService(storage, emailService, envConfig, log)

	validate = validator.New()

	resolver := app.NewResolver(
		authService,
		userService,
		validate,
		envConfig,
		log,
	)

	r := router.Setup(envConfig.Environment)
	server := app.NewServer(r, resolver)

	err = server.Run()

	if err != nil {
		return err
	}

	return nil
}
