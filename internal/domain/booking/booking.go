package booking

import (
	"context"
	"time"
)

type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusCancelled BookingStatus = "cancelled"
)

type Booking struct {
	ID          string        `json:"id" db:"id"`
	UserID      string        `json:"user_id" db:"user_id"`
	EventID     string        `json:"event_id" db:"event_id"`
	Quantity    int           `json:"quantity" db:"quantity"`
	TotalAmount float64       `json:"total_amount" db:"total_amount"`
	Status      BookingStatus `json:"status" db:"status"`
	BookingTime time.Time     `json:"booking_time" db:"booking_time"`
	CancelledAt *time.Time    `json:"cancelled_at,omitempty" db:"cancelled_at"`
	CreatedAt   time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at" db:"updated_at"`
}

type BookingRepository interface {
	Create(booking *Booking) error
	Update(booking *Booking) error
	GetByID(id string) (*Booking, error)
	GetByUserID(userID string, limit, offset int) ([]*Booking, error)
	GetByEventID(eventID string, limit, offset int) ([]*Booking, error)
	CountByEventID(eventID string) (int, error)
	GetTotalBookings() (int64, error)
	GetBookingAnalytics(eventID string) (*BookingAnalytics, error)
}

// Analytics models

type BookingAnalytics struct {
	EventID       string  `json:"event_id"`
	TotalBookings int     `json:"total_bookings"`
	TotalRevenue  float64 `json:"total_revenue"`
	Confirmed     int     `json:"confirmed"`
	Cancelled     int     `json:"cancelled"`
	Pending       int     `json:"pending"`
}

type BookingUsecase interface {
	CreateBooking(ctx context.Context, booking *Booking) error
	CancelBooking(ctx context.Context, bookingID, userID string) error
	GetBooking(ctx context.Context, bookingID string) (*Booking, error)
	GetUserBookings(ctx context.Context, userID string, limit, offset int) ([]*Booking, error)
	GetEventBookings(ctx context.Context, eventID string, limit, offset int) ([]*Booking, error)
	GetBookingAnalytics(ctx context.Context, eventID string) (*BookingAnalytics, error)
}
