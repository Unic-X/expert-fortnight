package routes

import (
	"evently/internal/delivery/http/handler"
	"evently/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

func SetupEventRoutes(router *gin.Engine, eventHandler *handler.EventHandler, jwtMiddleware *middleware.JWTConfig) {
	eventGroup := router.Group("/api/events")
	{
		// Public routes
		eventGroup.GET("", eventHandler.ListUpcomingEvents)
		eventGroup.GET("/:id", eventHandler.GetEvent)

		// Protected routes (admin only)
		adminGroup := eventGroup.Group("")
		adminGroup.Use(jwtMiddleware.AuthMiddleware())
		adminGroup.Use(middleware.AdminMiddleware())
		{
			adminGroup.POST("", eventHandler.CreateEvent)
			adminGroup.PUT("/:id", eventHandler.UpdateEvent)
			adminGroup.DELETE("/:id", eventHandler.DeleteEvent)
		}
	}
}
