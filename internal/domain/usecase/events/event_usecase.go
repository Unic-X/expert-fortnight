package events

import "evently/internal/domain/model"

type EventUsecase interface {
	BrowseEvents() ([]*model.Event, error)
	BookTicket(userID, eventID string) (*model.Booking, error)
	CancelBooking(userID, bookingID string) error
}
