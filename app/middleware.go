package app

import (
	"fmt"
	"github.com/VerzCar/vyf-lib-awsx"
	routerContext "github.com/VerzCar/vyf-user/app/router/ctx"
	"github.com/VerzCar/vyf-user/app/router/header"
	"github.com/gin-gonic/gin"
	"net/http"
)

// authGuard verifies the Authorization token against the SSO service.
// If the authentification fails the request will be aborted.
// Otherwise, the given subject of the token will be saved in the context and
// the next request served.
func (s *Server) authGuard(authService awsx.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		accessToken, err := header.Authorization(ctx, "Bearer")

		if err != nil {
			ctx.String(http.StatusUnauthorized, fmt.Sprintf("error: %s", err))
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		token, err := authService.DecodeAccessToken(ctx, accessToken)

		if err != nil {
			ctx.String(http.StatusUnauthorized, fmt.Sprintf("error decoding token"))
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		routerContext.SetAuthClaimsContext(ctx, token)
		routerContext.SetBearerTokenContext(ctx, accessToken)
		ctx.Next()
	}
}
