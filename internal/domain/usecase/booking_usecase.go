package usecase

import (
	"context"
	"evently/internal/domain/model"
)

type BookingUsecase interface {
	CreateBooking(ctx context.Context, booking *model.Booking) error
	CancelBooking(ctx context.Context, bookingID, userID string) error
	GetBooking(ctx context.Context, bookingID string) (*model.Booking, error)
	GetUserBookings(ctx context.Context, userID string, limit, offset int) ([]*model.Booking, error)
	GetEventBookings(ctx context.Context, eventID string, limit, offset int) ([]*model.Booking, error)
	GetBookingAnalytics(ctx context.Context, eventID string) (*model.BookingAnalytics, error)
}
