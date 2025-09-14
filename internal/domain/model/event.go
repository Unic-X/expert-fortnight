package model

import "time"

type Event struct {
	ID             string    `json:"id" db:"id"`
	Name           string    `json:"name" db:"name"`
	Description    string    `json:"description" db:"description"`
	Venue          string    `json:"venue" db:"venue"`
	EventTime      time.Time `json:"event_time" db:"event_time"`
	TotalCapacity  int       `json:"total_capacity" db:"total_capacity"`
	AvailableSeats int       `json:"available_seats" db:"available_seats"`
	Price          float64   `json:"price" db:"price"`
	CreatedBy      string    `json:"created_by" db:"created_by"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

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

type EventRepository interface {
	Create(event *Event) error
	Update(event *Event) error
	Delete(id string) error
	GetByID(id string) (*Event, error)
	ListUpcoming(limit, offset int) ([]*Event, error)
	ListAll(limit, offset int) ([]*Event, error)
	UpdateAvailableSeats(eventID string, quantity int) error
	GetMostPopularEvents(limit int) ([]*EventAnalytics, error)
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
type EventAnalytics struct {
	EventID         string  `json:"event_id" db:"event_id"`
	EventName       string  `json:"event_name" db:"event_name"`
	TotalBookings   int     `json:"total_bookings" db:"total_bookings"`
	TotalRevenue    float64 `json:"total_revenue" db:"total_revenue"`
	CapacityUsed    int     `json:"capacity_used" db:"capacity_used"`
	CapacityTotal   int     `json:"capacity_total" db:"capacity_total"`
	UtilizationRate float64 `json:"utilization_rate"`
}

type BookingAnalytics struct {
	EventID       string  `json:"event_id"`
	TotalBookings int     `json:"total_bookings"`
	TotalRevenue  float64 `json:"total_revenue"`
	Confirmed     int     `json:"confirmed"`
	Cancelled     int     `json:"cancelled"`
	Pending       int     `json:"pending"`
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
