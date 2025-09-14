package impl

import (
	"context"
	"fmt"
	"time"

	"evently/internal/domain/events"
	"evently/internal/domain/model"
	"evently/internal/domain/usecase"

	"github.com/google/uuid"
)

type notificationUsecaseImpl struct {
	notificationRepo model.NotificationRepository
	eventRepo        events.EventRepository
}

func NewNotificationUsecase(
	notificationRepo model.NotificationRepository,
	eventRepo events.EventRepository,
) usecase.NotificationUsecase {
	return &notificationUsecaseImpl{
		notificationRepo: notificationRepo,
		eventRepo:        eventRepo,
	}
}

func (u *notificationUsecaseImpl) CreateNotification(ctx context.Context, notification *model.Notification) error {
	// Generate ID if not provided
	if notification.ID == "" {
		notification.ID = uuid.New().String()
	}

	// Set timestamps
	now := time.Now()
	notification.CreatedAt = now
	notification.UpdatedAt = now

	return u.notificationRepo.Create(notification)
}

func (u *notificationUsecaseImpl) GetUserNotifications(ctx context.Context, userID string, limit, offset int) ([]*model.Notification, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	return u.notificationRepo.GetByUserID(userID, limit, offset)
}

func (u *notificationUsecaseImpl) MarkNotificationAsRead(ctx context.Context, notificationID string) error {
	return u.notificationRepo.MarkAsRead(notificationID)
}

func (u *notificationUsecaseImpl) MarkAllNotificationsAsRead(ctx context.Context, userID string) error {
	return u.notificationRepo.MarkAllAsRead(userID)
}

func (u *notificationUsecaseImpl) GetUnreadNotificationCount(ctx context.Context, userID string) (int, error) {
	return u.notificationRepo.GetUnreadCount(userID)
}

func (u *notificationUsecaseImpl) SendWaitlistNotification(ctx context.Context, userID, eventID string, quantity int) error {
	// Get event details
	event, err := u.eventRepo.GetByID(eventID)
	if err != nil {
		return fmt.Errorf("failed to get event details: %w", err)
	}

	// Create notification
	notification := &model.Notification{
		ID:        uuid.New().String(),
		UserID:    userID,
		EventID:   eventID,
		Type:      model.NotificationTypeWaitlistSpotAvailable,
		Title:     "Waitlist Spot Available",
		Message:   fmt.Sprintf("Great news! A spot is now available for %s. You can book %d ticket(s). This offer expires in 30 minutes.", event.Name, quantity),
		IsRead:    false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return u.notificationRepo.Create(notification)
}
