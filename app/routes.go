package app

func (s *Server) routes() {
	router := s.router

	// Authorization group
	authorized := router.Group("/")
	authorized.Use(s.ginContextToContext())
	authorized.Use(s.authGuard(s.resolver.ssoService))
	{
		// graphql route
		authorized.POST("/query", gqlHandler(s.resolver))
	}

}
