package api_test

import (
	"context"
	"gitlab.vecomentman.com/libs/awsx"
	"gitlab.vecomentman.com/libs/logger"
	appConfig "gitlab.vecomentman.com/vote-your-face/service/user/app/config"
	"gitlab.vecomentman.com/vote-your-face/service/user/utils"
	"os"
	"testing"
)

type MockUser struct {
	Elon    *awsx.JWTToken
	NewUser *awsx.JWTToken
}

var (
	config   *appConfig.Config
	log      logger.Logger
	ctx      context.Context
	mockUser *MockUser
)

func TestMain(m *testing.M) {
	configPath := utils.FromBase("app/config/")

	config = appConfig.NewConfig(configPath)
	log = logger.NewLogger(configPath)

	ctx = context.Background()

	mockUser = &MockUser{
		Elon:    &awsx.JWTToken{Subject: "elon"},
		NewUser: &awsx.JWTToken{Subject: "newUser"},
	}

	code := m.Run()

	os.Exit(code)
}

// login the user into the auth context
func putUserIntoContext(jwtToken *awsx.JWTToken) context.Context {
	ctx = context.WithValue(context.Background(), "AuthClaimsContextKey", jwtToken)
	return ctx
}

// emptyUserContext represents an authentication token in context without value
func emptyUserContext() context.Context {
	ctx = context.WithValue(context.Background(), "AuthClaimsContextKey", nil)
	return ctx
}
