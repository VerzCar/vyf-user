package ctx

import (
	"context"
	"fmt"
	"github.com/VerzCar/vyf-lib-awsx"
	"github.com/gin-gonic/gin"
)

type scopedContextKey string

const authClaimsContextKey = scopedContextKey("AuthClaimsContextKey")
const bearerTokenContextKey = scopedContextKey("BearerTokenContextKey")

func SetAuthClaimsContext(ctx *gin.Context, val interface{}) {
	c := context.WithValue(ctx.Request.Context(), authClaimsContextKey, val)
	ctx.Request = ctx.Request.WithContext(c)
}

func SetBearerTokenContext(ctx *gin.Context, val string) {
	c := context.WithValue(ctx.Request.Context(), bearerTokenContextKey, val)
	ctx.Request = ctx.Request.WithContext(c)
}

func ContextToAuthClaims(ctx context.Context) (*awsx.JWTToken, error) {
	authClaimsValue := ctx.Value(authClaimsContextKey)

	if authClaimsValue == nil {
		err := fmt.Errorf("could not retrieve auth claims")
		return nil, err
	}

	authClaims, ok := authClaimsValue.(*awsx.JWTToken)

	if !ok {
		err := fmt.Errorf("auth claims has wrong type")
		return nil, err
	}

	return authClaims, nil
}

func ContextToBearerToken(ctx context.Context) (string, error) {
	bearerTokenValue := ctx.Value(bearerTokenContextKey)

	if bearerTokenValue == nil {
		err := fmt.Errorf("could not retrieve bearer token")
		return "", err
	}

	bearerToken, ok := bearerTokenValue.(string)

	if !ok {
		err := fmt.Errorf("bearer token has wrong type")
		return "", err
	}

	return bearerToken, nil
}
