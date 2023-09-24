package app

func (s *Server) routes() {
	router := s.router

	// Service group
	v1 := router.Group("/v1/api/user")

	// Authorization group
	authorized := v1.Group("/")
	authorized.Use(s.authGuard(s.authService))
	{
		authorized.GET("/me", s.User())
		authorized.PUT("/update", s.UpdateUser())
	}
}
