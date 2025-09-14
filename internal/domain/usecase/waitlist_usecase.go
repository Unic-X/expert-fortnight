package usecase

import (
	"context"
	"evently/internal/domain/model"
)

type WaitlistUsecase interface {
	JoinWaitlist(ctx context.Context, userID, eventID string, quantity int) error
	LeaveWaitlist(ctx context.Context, userID, eventID string) error
	GetUserWaitlist(ctx context.Context, userID string, limit, offset int) ([]*model.Waitlist, error)
	GetEventWaitlist(ctx context.Context, eventID string, limit, offset int) ([]*model.Waitlist, error)
	ProcessWaitlistNotifications(ctx context.Context, eventID string, availableQuantity int) error
	GetWaitlistPosition(ctx context.Context, userID, eventID string) (int, error)
	ConvertWaitlistToBooking(ctx context.Context, waitlistID string) error
	CleanupExpiredWaitlist(ctx context.Context) error
	GetWaitlistByID(ctx context.Context, waitlistID string) (*model.Waitlist, error)
}

type NotificationUsecase interface {
	CreateNotification(ctx context.Context, notification *model.Notification) error
	GetUserNotifications(ctx context.Context, userID string, limit, offset int) ([]*model.Notification, error)
	MarkNotificationAsRead(ctx context.Context, notificationID string) error
	MarkAllNotificationsAsRead(ctx context.Context, userID string) error
	GetUnreadNotificationCount(ctx context.Context, userID string) (int, error)
	SendWaitlistNotification(ctx context.Context, userID, eventID string, quantity int) error
}
