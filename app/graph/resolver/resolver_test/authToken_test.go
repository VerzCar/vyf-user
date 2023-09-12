package resolver_test

/*
import (
	assertPkg "github.com/stretchr/testify/assert"
	"github.com/VerzCar/vyf-user/api/model"
	"github.com/VerzCar/vyf-user/app/cache/testcache"
	"github.com/VerzCar/vyf-user/app/database/testdb"
	"github.com/VerzCar/vyf-user/test/factory"
	"github.com/VerzCar/vyf-user/utils"
	"github.com/VerzCar/vyf-user/utils/testing/client"
	"testing"
	"time"
)

const (
	authToken = `mutation authToken($credentials: Credentials!) {
		authToken(credentials: $credentials) {
			token
			type
			expiresAt
		}
		}`
)

func TestAuthToken(t *testing.T) {
	assert := assertPkg.New(t)

	testdb.Reset(resolver.DB, resolver.Log)
	testcache.Reset(resolver.Rdb)
	factory.Setup(resolver.DB)

	var resp struct {
		AuthToken model.Token
	}

	err := c.Post(
		authToken,
		&resp,
		client.Var("credentials",
			model.Credentials{
				Email:    factory.Martin.Email,
				Password: factory.Martin.Password,
			}))

	assert.NoError(err)

	assert.Equal(resp.AuthToken.Type, resolver.Config.Token.Type)
	assert.NotEqual(resp.AuthToken.Token, "")
	expectedTime := time.Now()
	timeDelta := utils.FormatDuration(resolver.Config.Ttl.Token.Default)
	expectedTime.Add(timeDelta)
	assert.Greater(resp.AuthToken.ExpiresAt.String(), expectedTime.String())

}

func TestAuthToken_Assert_Wrong_Email(t *testing.T) {
	assert := assertPkg.New(t)

	testdb.Reset(resolver.DB, resolver.Log)
	testcache.Reset(resolver.Rdb)
	factory.Setup(resolver.DB)

	var resp struct {
		AuthToken model.Token
	}

	err := c.Post(
		authToken,
		&resp,
		client.Var("credentials",
			model.Credentials{
				Email:    "not.exist@email.com",
				Password: factory.Martin.Password,
			}))

	assert.Equal(err.Error(), `[{"message":"authentication failed","path":["authToken"]}]`)

}

func TestAuthToken_Assert_Wrong_Pwd(t *testing.T) {
	assert := assertPkg.New(t)

	testdb.Reset(resolver.DB, resolver.Log)
	testcache.Reset(resolver.Rdb)
	factory.Setup(resolver.DB)

	var resp struct {
		AuthToken model.Token
	}

	err := c.Post(
		authToken,
		&resp,
		client.Var("credentials",
			model.Credentials{
				Email:    factory.Martin.Email,
				Password: "wrongPassword123",
			}))

	assert.Equal(err.Error(), `[{"message":"authentication failed","path":["authToken"]}]`)

}

func TestAuthToken_Assert_Inactive_User(t *testing.T) {
	assert := assertPkg.New(t)

	testdb.Reset(resolver.DB, resolver.Log)
	testcache.Reset(resolver.Rdb)
	factory.Setup(resolver.DB)

	var resp struct {
		AuthToken model.Token
	}

	err := c.Post(
		authToken,
		&resp,
		client.Var("credentials",
			model.Credentials{
				Email:    factory.Albert.Email,
				Password: factory.Albert.Password,
			}))

	assert.Equal(err.Error(), `[{"message":"authentication failed","path":["authToken"]}]`)

}
*/
