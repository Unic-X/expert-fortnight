package routes

import "github.com/gin-gonic/gin"

// Create, update, and manage events.
// View booking analytics (total bookings, most popular events, capacity utilization).

func AdminRoutes(r gin.IRouter) {
	r.POST("/event")
	r.PUT("/event/:id")
	r.DELETE("/event/:id")
	r.GET("/analytics")
}
