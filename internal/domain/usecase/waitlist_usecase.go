package usecase

import (
	"context"
	"evently/internal/domain/model"
)

type NotificationUsecase interface {
	CreateNotification(ctx context.Context, notification *model.Notification) error
	GetUserNotifications(ctx context.Context, userID string, limit, offset int) ([]*model.Notification, error)
	MarkNotificationAsRead(ctx context.Context, notificationID string) error
	MarkAllNotificationsAsRead(ctx context.Context, userID string) error
	GetUnreadNotificationCount(ctx context.Context, userID string) (int, error)
	SendWaitlistNotification(ctx context.Context, userID, eventID string, quantity int) error
}
