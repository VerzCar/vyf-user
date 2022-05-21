package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.vecomentman.com/libs/sso"
	routerContext "gitlab.vecomentman.com/service/user/app/router/ctx"
	"gitlab.vecomentman.com/service/user/app/router/header"
	"net/http"
)

// ginContextToContext creates a gin middleware to add its context
// to the context.Context
func (s *Server) ginContextToContext() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		routerContext.SetGinContext(ctx)
		ctx.Next()
	}
}

// authGuard verifies the Authorization token against the SSO service.
// If the authentification fails the request will be aborted.
// Otherwise, the given subject of the token will be saved in the context and
// the next request served.
func (s *Server) authGuard(ssoService sso.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		accessToken, err := header.Authorization(ctx, "Bearer")

		if err != nil {
			ctx.String(http.StatusUnauthorized, fmt.Sprintf("error: %s", err))
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		token, ssoClaims, err := ssoService.DecodeAccessToken(ctx, accessToken)

		if err != nil {
			ctx.String(http.StatusUnauthorized, fmt.Sprintf("error decoding token"))
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			ctx.String(http.StatusUnauthorized, fmt.Sprint("token not valid"))
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		routerContext.SetSsoClaimsContext(ctx, ssoClaims)
		ctx.Next()
	}
}
