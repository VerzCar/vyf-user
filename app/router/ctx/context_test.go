package ctx

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"gitlab.vecomentman.com/libs/sso"
	"net/http"
	"testing"
)

func TestSetGinContext(t *testing.T) {
	req := &http.Request{}
	parentGinContext := &gin.Context{
		Request: req,
	}
	SetGinContext(parentGinContext)

	ginContext := parentGinContext.Request.Context().Value(ginContextKey)
	require.NotNil(t, ginContext)

	_, ok := ginContext.(*gin.Context)

	require.True(t, ok)
}

func TestContextToGinContext(t *testing.T) {
	ctx := context.Background()
	parentGinContext := &gin.Context{}
	testContext := context.WithValue(ctx, ginContextKey, parentGinContext)
	ginCtx, err := ContextToGinContext(testContext)

	require.Nil(t, err)
	require.Empty(t, ginCtx.Params.ByName("example"))
}

func TestSetSsoClaimsContext(t *testing.T) {
	req := &http.Request{}
	parentGinContext := &gin.Context{
		Request: req,
	}
	ssoValue := &sso.SsoClaims{
		Name: "Example",
	}
	SetSsoClaimsContext(parentGinContext, ssoValue)

	ssoContext := parentGinContext.Request.Context().Value(ssoContextKey)
	require.NotNil(t, ssoContext)

	contextSsoValue, ok := ssoContext.(*sso.SsoClaims)

	require.True(t, ok)
	require.Equal(t, contextSsoValue, ssoValue)
}

func TestContextToSsoClaims(t *testing.T) {
	ctx := context.Background()
	ssoValue := &sso.SsoClaims{
		Name: "Example",
	}

	testContext := context.WithValue(ctx, ssoContextKey, ssoValue)
	contextSsoValue, err := ContextToSsoClaims(testContext)

	require.Nil(t, err)
	require.Equal(t, contextSsoValue, ssoValue)
}
