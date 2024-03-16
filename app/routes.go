package app

import cors "github.com/rs/cors/wrapper/gin"

func (s *Server) routes() {
	router := s.router

	// Service group
	v1 := router.Group("/v1/api/user", cors.AllowAll())

	// Authorization group
	authorized := v1.Group("/")
	authorized.Use(s.authGuard(s.authService))
	authorized.Use(cors.AllowAll())
	{
		authorized.GET("", s.UserMe())
		authorized.GET("/:identityId", s.UserX())
		authorized.GET("/users", s.Users())
		authorized.GET("/users/:username", s.UsersByUsername())

		authorized.PUT("/update", s.UpdateUser())

		// Upload group
		upload := authorized.Group("/upload")
		upload.PUT("/profile-img", s.UploadProfileImage())
	}
}
