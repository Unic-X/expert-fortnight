package model

import "time"

type Event struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Venue     string    `json:"venue"`
	Time      time.Time `json:"time"`
	Capacity  int       `json:"capacity"`
	CreatedAt time.Time `json:"created_at"`
}

type Booking struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	EventID   string    `json:"event_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type EventRepository interface {
	Create(event *Event) error
	Update(event *Event) error
	Delete(id string) error
	GetByID(id string) (*Event, error)
	ListUpcoming() ([]*Event, error)
	DecrementCapacity(eventID string) error
	IncrementCapacity(eventID string) error
}

type BookingRepository interface {
	Create(booking *Booking) error
	Cancel(bookingID string) error
	GetByID(id string) (*Booking, error)
	GetByUser(userID string) ([]*Booking, error)
	CountByEvent(eventID string) (int, error)
}
