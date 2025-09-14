package waitlist

import (
	"context"
	"time"
)

type Waitlist struct {
	ID         string         `json:"id" db:"id"`
	UserID     string         `json:"user_id" db:"user_id"`
	EventID    string         `json:"event_id" db:"event_id"`
	Quantity   int            `json:"quantity" db:"quantity"`
	Priority   int            `json:"priority" db:"priority"`
	Status     WaitlistStatus `json:"status" db:"status"`
	JoinedAt   time.Time      `json:"joined_at" db:"joined_at"`
	NotifiedAt *time.Time     `json:"notified_at,omitempty" db:"notified_at"`
	ExpiresAt  *time.Time     `json:"expires_at,omitempty" db:"expires_at"`
	CreatedAt  time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at" db:"updated_at"`
}

// Waitlist models
type WaitlistStatus string

const (
	WaitlistStatusActive    WaitlistStatus = "active"
	WaitlistStatusNotified  WaitlistStatus = "notified"
	WaitlistStatusExpired   WaitlistStatus = "expired"
	WaitlistStatusConverted WaitlistStatus = "converted"
)

type WaitlistRepository interface {
	Create(waitlist *Waitlist) error
	Update(waitlist *Waitlist) error
	Delete(id string) error
	GetByID(id string) (*Waitlist, error)
	GetByUserAndEvent(userID, eventID string) (*Waitlist, error)
	GetByEventID(eventID string, limit, offset int) ([]*Waitlist, error)
	GetByUserID(userID string, limit, offset int) ([]*Waitlist, error)
	GetNextInQueue(eventID string, quantity int) ([]*Waitlist, error)
	CountByEventID(eventID string) (int, error)
	UpdateStatus(id string, status WaitlistStatus) error
	CleanupExpired() error
}

type WaitlistUsecase interface {
	JoinWaitlist(ctx context.Context, userID, eventID string, quantity int) error
	LeaveWaitlist(ctx context.Context, userID, eventID string) error
	GetUserWaitlist(ctx context.Context, userID string, limit, offset int) ([]*Waitlist, error)
	GetEventWaitlist(ctx context.Context, eventID string, limit, offset int) ([]*Waitlist, error)
	ProcessWaitlistNotifications(ctx context.Context, eventID string, availableQuantity int) error
	GetWaitlistPosition(ctx context.Context, userID, eventID string) (int, error)
	ConvertWaitlistToBooking(ctx context.Context, waitlistID string) error
	CleanupExpiredWaitlist(ctx context.Context) error
	GetWaitlistByID(ctx context.Context, waitlistID string) (*Waitlist, error)
}
