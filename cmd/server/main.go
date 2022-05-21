package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"gitlab.vecomentman.com/libs/logger"
	"gitlab.vecomentman.com/libs/sso"
	"gitlab.vecomentman.com/service/user/api"
	"gitlab.vecomentman.com/service/user/app"
	"gitlab.vecomentman.com/service/user/app/cache"
	"gitlab.vecomentman.com/service/user/app/config"
	"gitlab.vecomentman.com/service/user/app/database"
	"gitlab.vecomentman.com/service/user/app/email"
	"gitlab.vecomentman.com/service/user/app/graph/client"
	"gitlab.vecomentman.com/service/user/app/router"
	"gitlab.vecomentman.com/service/user/repository"
	"gitlab.vecomentman.com/service/user/services"
	"gitlab.vecomentman.com/service/user/utils"
	"net/http"
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

	redisCache := cache.Connect(log, envConfig)

	storage := repository.NewStorage(db, envConfig, log)

	sqlDb, _ := db.DB()
	err := storage.RunMigrationsUp(sqlDb)

	if err != nil {
		return err
	}

	redis := cache.NewRedisCache(redisCache, envConfig, log)

	// initialize auth service
	ssoService := sso.NewService(
		envConfig.Hosts.Svc.Sso,
		sso.Realm(envConfig.Sso.Realm.Name),
		sso.ClientId(envConfig.Sso.Realm.Client.Id),
		sso.ClientSecret(envConfig.Sso.Realm.Client.Secret),
		sso.AdminCredentials(
			envConfig.Sso.Realm.Admin.Username,
			envConfig.Sso.Realm.Admin.Password,
		),
	)

	// initialize third party services
	httpClient := &http.Client{}
	gqlClient := client.New(httpClient)
	paymentSvc := services.NewPaymentService(envConfig.Hosts.Svc.Payment, log, gqlClient)
	emailService := email.NewService(envConfig, log)

	if err != nil {
		return err
	}

	// initialize api services
	userService := api.NewUserService(storage, redis, ssoService, emailService, envConfig, log)
	companyService := api.NewCompanyService(storage, redis, emailService, paymentSvc, envConfig, log)

	validate = validator.New()

	resolver := app.NewResolver(
		ssoService,
		userService,
		companyService,
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
