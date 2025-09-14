package model

import "time"

// Notification models
type NotificationType string

const (
	NotificationTypeWaitlistSpotAvailable NotificationType = "waitlist_spot_available"
	NotificationTypeBookingConfirmed      NotificationType = "booking_confirmed"
	NotificationTypeBookingCancelled      NotificationType = "booking_cancelled"
)

type Notification struct {
	ID        string           `json:"id" db:"id"`
	UserID    string           `json:"user_id" db:"user_id"`
	EventID   string           `json:"event_id" db:"event_id"`
	Type      NotificationType `json:"type" db:"type"`
	Title     string           `json:"title" db:"title"`
	Message   string           `json:"message" db:"message"`
	IsRead    bool             `json:"is_read" db:"is_read"`
	CreatedAt time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt time.Time        `json:"updated_at" db:"updated_at"`
}

type NotificationRepository interface {
	Create(notification *Notification) error
	GetByUserID(userID string, limit, offset int) ([]*Notification, error)
	MarkAsRead(id string) error
	MarkAllAsRead(userID string) error
	GetUnreadCount(userID string) (int, error)
}
