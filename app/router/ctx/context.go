package ctx

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.vecomentman.com/libs/sso"
)

const ginContextKey = "GinCtxKey"
const ssoContextKey = "SsoCtxKey"

func SetGinContext(ctx *gin.Context) {
	c := context.WithValue(ctx.Request.Context(), ginContextKey, ctx)
	ctx.Request = ctx.Request.WithContext(c)
}

func ContextToGinContext(ctx context.Context) (*gin.Context, error) {
	ginContext := ctx.Value(ginContextKey)

	if ginContext == nil {
		err := fmt.Errorf("could not retrieve gin.Context")
		return nil, err
	}

	gc, ok := ginContext.(*gin.Context)

	if !ok {
		err := fmt.Errorf("gin.Context has wrong type")
		return nil, err
	}

	return gc, nil
}

func SetSsoClaimsContext(ctx *gin.Context, val interface{}) {
	c := context.WithValue(ctx.Request.Context(), ssoContextKey, val)
	ctx.Request = ctx.Request.WithContext(c)
}

func ContextToSsoClaims(ctx context.Context) (*sso.Claims, error) {
	ssoClaimsValue := ctx.Value(ssoContextKey)

	if ssoClaimsValue == nil {
		err := fmt.Errorf("could not retrieve sso claims")
		return nil, err
	}

	ssoClaims, ok := ssoClaimsValue.(*sso.Claims)

	if !ok {
		err := fmt.Errorf("sso claims has wrong type")
		return nil, err
	}

	return ssoClaims, nil
}
