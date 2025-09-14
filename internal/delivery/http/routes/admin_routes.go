package routes

import (
	"evently/internal/delivery/http/handler"
	"evently/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

// Create, update, and manage events.
// View booking analytics (total bookings, most popular events, capacity utilization).

func SetupAdminRoutes(router *gin.Engine, adminHandler *handler.AdminHandler, jwtMiddleware *middleware.JWTConfig) {
	adminGroup := router.Group("/api/admin")
	adminGroup.Use(jwtMiddleware.AuthMiddleware())
	adminGroup.Use(middleware.AdminMiddleware())
	{
		adminGroup.GET("/events", adminHandler.GetAllEvents)
		adminGroup.GET("/events/:eventId/bookings", adminHandler.GetEventBookings)
		adminGroup.GET("/events/:eventId/analytics", adminHandler.GetBookingAnalytics)
		adminGroup.GET("/analytics/events", adminHandler.GetEventAnalytics)
	}
}
