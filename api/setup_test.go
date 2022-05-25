package api_test

import (
	"context"
	"gitlab.vecomentman.com/libs/logger"
	"gitlab.vecomentman.com/libs/sso"
	appConfig "gitlab.vecomentman.com/vote-your-face/service/user/app/config"
	"gitlab.vecomentman.com/vote-your-face/service/user/test/factory"
	"gitlab.vecomentman.com/vote-your-face/service/user/utils"
	"os"
	"testing"
)

var (
	config *appConfig.Config
	log    logger.Logger
	dbUser factory.User
	ctx    context.Context
)

// Setup test env case
func TestMain(m *testing.M) {
	configPath := utils.FromBase("app/config/")

	config = appConfig.NewConfig(configPath)
	log = logger.NewLogger(configPath)

	dbUser = factory.NewUsers()

	ctx = context.Background()

	code := m.Run()

	os.Exit(code)
}

// login the user into the context
func ssoUserIntoContext(ssoUser *sso.SsoClaims) context.Context {
	ctx = context.WithValue(context.Background(), "SsoCtxKey", ssoUser)
	return ctx
}

// Let the context have an invalid context for the user
func ssoUserIntoContextFail(ssoUser interface{}) context.Context {
	ctx = context.WithValue(context.Background(), "failUnicorn555", ssoUser)
	return ctx
}
