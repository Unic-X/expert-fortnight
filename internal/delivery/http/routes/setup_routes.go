package routes

import (
	"evently/internal/delivery/http/handler"
	"evently/internal/delivery/http/middleware"
	"evently/internal/di"

	"github.com/gin-gonic/gin"
)

func AllRoutes(router *gin.Engine, container *di.Container, jwtMiddleware *middleware.JWTConfig) {

	authHandler := handler.NewAuthHandler(container.AuthUseCase)
	eventHandler := handler.NewEventHandler(container.EventUseCase)
	bookingHandler := handler.NewBookingHandler(container.BookingUseCase, container.WaitlistUseCase)
	adminHandler := handler.NewAdminHandler(container.EventUseCase, container.BookingUseCase)

	api := router.Group("/api")
	{
		SetupAuthRoutes(api, authHandler)
		SetupEventRoutes(api, eventHandler, jwtMiddleware)
		SetupBookingRoutes(api, bookingHandler, jwtMiddleware)
		SetupAdminRoutes(api, adminHandler, jwtMiddleware)
	}
}
