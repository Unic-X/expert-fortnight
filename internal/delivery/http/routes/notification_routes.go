package routes

import (
	"evently/internal/delivery/http/handler"
	"evently/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

func SetupNotificationRoutes(router *gin.Engine, notificationHandler *handler.NotificationHandler, jwtMiddleware *middleware.JWTConfig) {
	notificationGroup := router.Group("/api/notifications")
	notificationGroup.Use(jwtMiddleware.AuthMiddleware())
	{
		notificationGroup.GET("", notificationHandler.GetUserNotifications)
		notificationGroup.PUT("/:id/read", notificationHandler.MarkNotificationAsRead)
		notificationGroup.PUT("/read-all", notificationHandler.MarkAllNotificationsAsRead)
		notificationGroup.GET("/unread-count", notificationHandler.GetUnreadNotificationCount)
	}
}
