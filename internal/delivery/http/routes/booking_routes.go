package routes

import (
	"evently/internal/delivery/http/handler"
	"evently/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

func SetupBookingRoutes(router *gin.Engine, bookingHandler *handler.BookingHandler, jwtMiddleware *middleware.JWTConfig) {
	bookingGroup := router.Group("/api/bookings")
	bookingGroup.Use(jwtMiddleware.AuthMiddleware())
	{
		bookingGroup.POST("", bookingHandler.CreateBooking)
		bookingGroup.GET("/my", bookingHandler.GetUserBookings)
		bookingGroup.GET("/:id", bookingHandler.GetBooking)
		bookingGroup.PUT("/:id/cancel", bookingHandler.CancelBooking)
	}
}
