package server

import (
	"os"

	"evently/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router    *gin.Engine
	jwtConfig *middleware.JWTConfig
}

func NewServer() *Server {
	r := gin.Default()
	jwtConfig := middleware.NewJWTConfig()

	srv := &Server{
		router:    r,
		jwtConfig: jwtConfig,
	}

	srv.SetupRoutes()
	return srv
}

func (s *Server) SetupRoutes() {
	protected := s.router.Group("/api")
	protected.Use(s.jwtConfig.AuthMiddleware())
	{
		protected.GET("/protected", func(c *gin.Context) {
			userID, userType, _ := middleware.GetUserFromContext(c)
			c.JSON(200, gin.H{
				"message":   "Access granted",
				"user_id":   userID,
				"user_type": userType,
			})
		})

	}
}

func (s *Server) Start() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return s.router.Run(":" + port)
}
