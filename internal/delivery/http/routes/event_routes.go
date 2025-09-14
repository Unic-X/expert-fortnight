package routes

import (
	"evently/internal/delivery/http/handler"
	"evently/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

func SetupEventRoutes(router *gin.RouterGroup, eventHandler *handler.EventHandler, jwtMiddleware *middleware.JWTConfig) {
	eventGroup := router.Group("/events")
	{
		// Public routes
		eventGroup.GET("", eventHandler.ListUpcomingEvents)
		eventGroup.GET("/:id", eventHandler.GetEvent)

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
