package main

import (
	"fmt"
	"github.com/VerzCar/vyf-lib-awsx"
	"github.com/VerzCar/vyf-lib-logger"
	"github.com/VerzCar/vyf-user/api"
	"github.com/VerzCar/vyf-user/app"
	"github.com/VerzCar/vyf-user/app/config"
	"github.com/VerzCar/vyf-user/app/database"
	"github.com/VerzCar/vyf-user/app/router"
	"github.com/VerzCar/vyf-user/repository"
	"github.com/VerzCar/vyf-user/utils"
	"github.com/go-playground/validator/v10"
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

	// initialize api services
	userService := api.NewUserService(storage, envConfig, log)

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
