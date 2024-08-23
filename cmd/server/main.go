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
		_, err := fmt.Fprintf(os.Stderr, "startup error: %s\\n", err)
		if err != nil {
			return
		}
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

	// initialize aws services
	authService, err := awsx.NewAuthService(
		awsx.AppClientId(envConfig.Aws.Auth.ClientId),
		awsx.ClientSecret(envConfig.Aws.Auth.ClientSecret),
		awsx.AwsDefaultRegion(envConfig.Aws.Auth.AwsDefaultRegion),
		awsx.UserPoolId(envConfig.Aws.Auth.UserPoolId),
	)

	if err != nil {
		return err
	}

	// initialize aws services
	s3Service, err := awsx.NewS3Service(
		awsx.AccessKeyID(envConfig.Aws.S3.AccessKeyId),
		awsx.AccessKeySecret(envConfig.Aws.S3.AccessKeySecret),
		awsx.Region(envConfig.Aws.S3.Region),
		awsx.BucketName(envConfig.Aws.S3.BucketName),
		awsx.DefaultBaseURL(envConfig.Aws.S3.DefaultBaseURL),
		awsx.UploadTimeout(envConfig.Aws.S3.UploadTimeout),
	)

	if err != nil {
		return err
	}

	// initialize api services
	userService := api.NewUserService(storage, envConfig, log)
	userUploadService := api.NewUserUploadService(userService, s3Service, envConfig, log)

	validate = validator.New()

	r := router.Setup(envConfig)
	server := app.NewServer(
		r,
		authService,
		userService,
		userUploadService,
		validate,
		envConfig,
		log,
	)

	err = server.Run()

	if err != nil {
		return err
	}

	return nil
}
