package ctx

import (
	"context"
	"github.com/VerzCar/vyf-lib-awsx"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestSetAuthClaimsContext(t *testing.T) {
	tests := []struct {
		name string
		val  interface{}
	}{
		{
			name: "Test Set Auth Claims Context",
			val:  &awsx.JWTToken{},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				ctx := &gin.Context{
					Request: &http.Request{},
				}
				SetAuthClaimsContext(ctx, tt.val)
				assert.NotNil(t, ctx.Request.Context().Value(authClaimsContextKey))
			},
		)
	}
}

func TestSetBearerTokenContext(t *testing.T) {
	tests := []struct {
		name string
		val  string
	}{
		{
			name: "Test Set Bearer Token Context",
			val:  "testToken",
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				ctx := &gin.Context{
					Request: &http.Request{},
				}
				SetBearerTokenContext(ctx, tt.val)
				assert.NotNil(t, ctx.Request.Context().Value(bearerTokenContextKey))
			},
		)
	}
}

func TestContextToAuthClaims(t *testing.T) {
	tests := []struct {
		name string
		val  interface{}
		err  bool
	}{
		{
			name: "Test Context To Auth Claims - Success",
			val:  &awsx.JWTToken{},
			err:  false,
		},
		{
			name: "Test Context To Auth Claims - Failure",
			val:  "invalid",
			err:  true,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				ctx := context.WithValue(context.Background(), authClaimsContextKey, tt.val)
				_, err := ContextToAuthClaims(ctx)
				if tt.err {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			},
		)
	}
}

func TestContextToBearerToken(t *testing.T) {
	tests := []struct {
		name string
		val  interface{}
		err  bool
	}{
		{
			name: "Test Context To Bearer Token - Success",
			val:  "testToken",
			err:  false,
		},
		{
			name: "Test Context To Bearer Token - Failure",
			val:  123,
			err:  true,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				ctx := context.WithValue(context.Background(), bearerTokenContextKey, tt.val)
				_, err := ContextToBearerToken(ctx)
				if tt.err {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			},
		)
	}
}
