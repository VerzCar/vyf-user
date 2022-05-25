package email_test

import (
	"gitlab.vecomentman.com/libs/logger"
	appConfig "gitlab.vecomentman.com/vote-your-face/service/user/app/config"
	"gitlab.vecomentman.com/vote-your-face/service/user/utils"
	"os"
	"testing"
)

var (
	config *appConfig.Config
	log    logger.Logger
)

// Setup test env case
func TestMain(m *testing.M) {
	configPath := utils.FromBase("app/config/")

	config = appConfig.NewConfig(configPath)
	log = logger.NewLogger(configPath)

	code := m.Run()

	os.Exit(code)
}
