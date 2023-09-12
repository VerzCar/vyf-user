package resolver_test

/*import (
	"github.com/gin-gonic/gin"
	"github.com/VerzCar/vyf-user/app/cache"
	testdb2 "github.com/VerzCar/vyf-user/app/database/testdb"
	gqlResolver "github.com/VerzCar/vyf-user/app/graph/resolver"
	mainRouter "github.com/VerzCar/vyf-user/app/router"
	"github.com/VerzCar/vyf-user/config"
	"github.com/VerzCar/vyf-lib-logger"
	"github.com/VerzCar/vyf-user/utils"
	"github.com/VerzCar/vyf-user/utils/testing/client"
	"os"
	"testing"
)

var (
	resolver gqlResolver.Resolver
	router   *gin.Engine
	c        *client.Client
)

// Setup test env case
func TestMain(m *testing.M) {
	configPath := utils.FromBase("config/")

	resolver.Config = config.Load(configPath)
	resolver.Log = logger.NewLogger(configPath)

	resolver.Rdb = cache.Connect(resolver.Log, resolver.Config)

	resolver.DB = testdb2.Connect(resolver.Log, resolver.Config)
	testdb2.Setup(resolver.DB, resolver.Log, resolver.Config)

	router = mainRouter.Setup(&resolver)

	c = client.New(router)

	code := m.Run()

	os.Exit(code)
}*/
