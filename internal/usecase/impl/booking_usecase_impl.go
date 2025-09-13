package impl

import (
	"context"
	"fmt"
	"sync"
	"time"

	"evently/internal/domain/model"
	"evently/internal/domain/usecase"

	"github.com/google/uuid"
)

type bookingUsecaseImpl struct {
	bookingRepo model.BookingRepository
	eventRepo   model.EventRepository
	mu          sync.RWMutex // For handling concurrent bookings
}

func NewBookingUsecase(bookingRepo model.BookingRepository, eventRepo model.EventRepository) usecase.BookingUsecase {
	return &bookingUsecaseImpl{
		bookingRepo: bookingRepo,
		eventRepo:   eventRepo,
	}
}

func (u *bookingUsecaseImpl) CreateBooking(ctx context.Context, booking *model.Booking) error {
	// Use mutex to handle concurrent bookings safely
	u.mu.Lock()
	defer u.mu.Unlock()
	
	// Get event details
	event, err := u.eventRepo.GetByID(booking.EventID)
	if err != nil {
		return fmt.Errorf("event not found: %w", err)
	}
	
	// Check if event is in the future
	if event.EventTime.Before(time.Now()) {
		return fmt.Errorf("cannot book tickets for past events")
	}
	
	// Check seat availability
	if event.AvailableSeats < booking.Quantity {
		return fmt.Errorf("insufficient seats available. Available: %d, Requested: %d", 
			event.AvailableSeats, booking.Quantity)
	}
	
	// Generate booking ID and set timestamps
	booking.ID = uuid.New().String()
	booking.Status = model.BookingStatusPending
	booking.BookingTime = time.Now()
	booking.CreatedAt = time.Now()
	booking.UpdatedAt = time.Now()
	booking.TotalAmount = float64(booking.Quantity) * event.Price
	
	// Validate booking
	if err := u.validateBooking(booking); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	
	// Create booking
	if err := u.bookingRepo.Create(booking); err != nil {
		return fmt.Errorf("failed to create booking: %w", err)
	}
	
	// Update available seats
	if err := u.eventRepo.UpdateAvailableSeats(booking.EventID, -booking.Quantity); err != nil {
		// TODO: Implement compensation logic here in production
		return fmt.Errorf("failed to update seat availability: %w", err)
	}
	
	// Update booking status to confirmed
	booking.Status = model.BookingStatusConfirmed
	booking.UpdatedAt = time.Now()
	
	return u.bookingRepo.Update(booking)
}

func (u *bookingUsecaseImpl) CancelBooking(ctx context.Context, bookingID, userID string) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	
	// Get booking
	booking, err := u.bookingRepo.GetByID(bookingID)
	if err != nil {
		return fmt.Errorf("booking not found: %w", err)
	}
	
	// Check if user owns the booking
	if booking.UserID != userID {
		return fmt.Errorf("unauthorized: booking belongs to different user")
	}
	
	// Check if booking can be cancelled
	if booking.Status == model.BookingStatusCancelled {
		return fmt.Errorf("booking is already cancelled")
	}
	
	// Get event to check cancellation policy (e.g., can't cancel within 24 hours)
	event, err := u.eventRepo.GetByID(booking.EventID)
	if err != nil {
		return fmt.Errorf("event not found: %w", err)
	}
	
	// Check if event has already passed
	if event.EventTime.Before(time.Now()) {
		return fmt.Errorf("cannot cancel booking for past events")
	}
	
	// Update booking status
	now := time.Now()
	booking.Status = model.BookingStatusCancelled
	booking.CancelledAt = &now
	booking.UpdatedAt = now
	
	// Update booking
	if err := u.bookingRepo.Update(booking); err != nil {
		return fmt.Errorf("failed to cancel booking: %w", err)
	}
	
	// Return seats to available pool
	if err := u.eventRepo.UpdateAvailableSeats(booking.EventID, booking.Quantity); err != nil {
		// TODO: Implement compensation logic here in production
		return fmt.Errorf("failed to update seat availability: %w", err)
	}
	
	return nil
}

func (u *bookingUsecaseImpl) GetBooking(ctx context.Context, bookingID string) (*model.Booking, error) {
	return u.bookingRepo.GetByID(bookingID)
}

func (u *bookingUsecaseImpl) GetUserBookings(ctx context.Context, userID string, limit, offset int) ([]*model.Booking, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	
	return u.bookingRepo.GetByUserID(userID, limit, offset)
}

func (u *bookingUsecaseImpl) GetEventBookings(ctx context.Context, eventID string, limit, offset int) ([]*model.Booking, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	
	return u.bookingRepo.GetByEventID(eventID, limit, offset)
}

func (u *bookingUsecaseImpl) GetBookingAnalytics(ctx context.Context, eventID string) (*model.BookingAnalytics, error) {
	return u.bookingRepo.GetBookingAnalytics(eventID)
}

func (u *bookingUsecaseImpl) validateBooking(booking *model.Booking) error {
	if booking.UserID == "" {
		return fmt.Errorf("user ID is required")
	}
	
	if booking.EventID == "" {
		return fmt.Errorf("event ID is required")
	}
	
	if booking.Quantity <= 0 {
		return fmt.Errorf("quantity must be positive")
	}
	
	if booking.Quantity > 10 { // Business rule: max 10 tickets per booking
		return fmt.Errorf("cannot book more than 10 tickets at once")
	}
	
	return nil
}
