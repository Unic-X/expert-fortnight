package routes

import "github.com/gin-gonic/gin"

// User Features
// Browse a list of upcoming events with details (name, venue, time, capacity).
// Book and cancel tickets, ensuring seat availability is updated correctly.
// View booking history.

func UserRoutes(r gin.IRouter) {
	r.POST("/register")
	r.POST("/login")
	r.GET("/history")
	r.POST("/book")
	r.POST("/cancel")
}
