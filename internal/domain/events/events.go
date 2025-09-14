package events

import (
	"context"
	"time"
)

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

type EventUsecase interface {
	CreateEvent(ctx context.Context, event *Event) error
	UpdateEvent(ctx context.Context, event *Event) error
	DeleteEvent(ctx context.Context, eventID string) error
	GetEvent(ctx context.Context, eventID string) (*Event, error)
	ListUpcomingEvents(ctx context.Context, limit, offset int) ([]*Event, error)
	ListAllEvents(ctx context.Context, limit, offset int) ([]*Event, error)
	GetMostPopularEvents(ctx context.Context, limit int) ([]*EventAnalytics, error)
}

type EventRepository interface {
	Create(event *Event) error
	Update(event *Event) error
	Delete(id string) error
	GetByID(id string) (*Event, error)
	ListUpcoming(limit, offset int) ([]*Event, error)
	ListAll(limit, offset int) ([]*Event, error)
	UpdateAvailableSeats(eventID string, quantity int) error
	GetMostPopularEvents(ctx context.Context, limit int) ([]*EventAnalytics, error)
}

type EventAnalytics struct {
	EventID         string  `json:"event_id" db:"event_id"`
	EventName       string  `json:"event_name" db:"event_name"`
	TotalBookings   int     `json:"total_bookings" db:"total_bookings"`
	TotalRevenue    float64 `json:"total_revenue" db:"total_revenue"`
	CapacityUsed    int     `json:"capacity_used" db:"capacity_used"`
	CapacityTotal   int     `json:"capacity_total" db:"capacity_total"`
	UtilizationRate float64 `json:"utilization_rate"`
}
