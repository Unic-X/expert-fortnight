package impl

import (
	"context"
	"fmt"
	"sync"
	"time"

	"evently/internal/domain/booking"
	"evently/internal/domain/events"
	"evently/internal/domain/waitlist"

	"github.com/google/uuid"
)

type bookingUsecaseImpl struct {
	bookingRepo     booking.BookingRepository
	eventRepo       events.EventRepository
	waitlistUsecase waitlist.WaitlistUsecase
	mu              sync.RWMutex // For handling concurrent bookings
}

func NewBookingUsecase(
	bookingRepo booking.BookingRepository,
	eventRepo events.EventRepository,
	waitlistUsecase waitlist.WaitlistUsecase,
) booking.BookingUsecase {
	return &bookingUsecaseImpl{
		bookingRepo:     bookingRepo,
		eventRepo:       eventRepo,
		waitlistUsecase: waitlistUsecase,
	}
}

func (u *bookingUsecaseImpl) CreateBooking(ctx context.Context, newBooking *booking.Booking) error {
	// Use mutex to handle concurrent bookings safely
	u.mu.Lock()
	defer u.mu.Unlock()

	// Get event details
	event, err := u.eventRepo.GetByID(newBooking.EventID)
	if err != nil {
		return fmt.Errorf("event not found: %w", err)
	}

	// Check if event is in the future
	if event.EventTime.Before(time.Now()) {
		return fmt.Errorf("cannot book tickets for past events")
	}

	// Check seat availability
	if event.AvailableSeats < newBooking.Quantity {
		return fmt.Errorf("insufficient seats available. Available: %d, Requested: %d",
			event.AvailableSeats, newBooking.Quantity)
	}

	// Generate booking ID and set timestamps
	newBooking.ID = uuid.New().String()
	newBooking.Status = booking.BookingStatusPending
	newBooking.BookingTime = time.Now()
	newBooking.CreatedAt = time.Now()
	newBooking.UpdatedAt = time.Now()
	newBooking.TotalAmount = float64(newBooking.Quantity) * event.Price

	// Validate booking
	if err := u.validateBooking(newBooking); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Create booking
	if err := u.bookingRepo.Create(newBooking); err != nil {
		return fmt.Errorf("failed to create booking: %w", err)
	}

	// Update available seats
	if err := u.eventRepo.UpdateAvailableSeats(newBooking.EventID, -newBooking.Quantity); err != nil {
		// TODO: Implement compensation logic here in production
		return fmt.Errorf("failed to update seat availability: %w", err)
	}

	// Update booking status to confirmed
	newBooking.Status = booking.BookingStatusConfirmed
	newBooking.UpdatedAt = time.Now()

	return u.bookingRepo.Update(newBooking)
}

func (u *bookingUsecaseImpl) CancelBooking(ctx context.Context, bookingID, userID string) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	// Get booking
	oldBooking, err := u.bookingRepo.GetByID(bookingID)
	if err != nil {
		return fmt.Errorf("booking not found: %w", err)
	}

	// Check if user owns the booking
	if oldBooking.UserID != userID {
		return fmt.Errorf("unauthorized: booking belongs to different user")
	}

	// Check if booking can be cancelled
	if oldBooking.Status == booking.BookingStatusCancelled {
		return fmt.Errorf("booking is already cancelled")
	}

	// Get event to check cancellation policy (e.g., can't cancel within 24 hours)
	event, err := u.eventRepo.GetByID(oldBooking.EventID)
	if err != nil {
		return fmt.Errorf("event not found: %w", err)
	}

	// Check if event has already passed
	if event.EventTime.Before(time.Now()) {
		return fmt.Errorf("cannot cancel booking for past events")
	}

	// Update booking status
	now := time.Now()
	oldBooking.Status = booking.BookingStatusCancelled
	oldBooking.CancelledAt = &now
	oldBooking.UpdatedAt = now

	// Update booking
	if err := u.bookingRepo.Update(oldBooking); err != nil {
		return fmt.Errorf("failed to cancel booking: %w", err)
	}

	// Return seats to available pool
	if err := u.eventRepo.UpdateAvailableSeats(oldBooking.EventID, oldBooking.Quantity); err != nil {
		// TODO: Implement compensation logic here in production
		return fmt.Errorf("failed to update seat availability: %w", err)
	}

	// Process waitlist notifications for newly available seats
	if u.waitlistUsecase != nil {
		if err := u.waitlistUsecase.ProcessWaitlistNotifications(ctx, oldBooking.EventID, oldBooking.Quantity); err != nil {
			// Log error but don't fail the cancellation
			fmt.Printf("Failed to process waitlist notifications: %v\n", err)
		}
	}

	return nil
}

func (u *bookingUsecaseImpl) GetBooking(ctx context.Context, bookingID string) (*booking.Booking, error) {
	return u.bookingRepo.GetByID(bookingID)
}

func (u *bookingUsecaseImpl) GetUserBookings(ctx context.Context, userID string, limit, offset int) ([]*booking.Booking, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	return u.bookingRepo.GetByUserID(userID, limit, offset)
}

func (u *bookingUsecaseImpl) GetEventBookings(ctx context.Context, eventID string, limit, offset int) ([]*booking.Booking, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	return u.bookingRepo.GetByEventID(eventID, limit, offset)
}

func (u *bookingUsecaseImpl) GetBookingAnalytics(ctx context.Context, eventID string) (*booking.BookingAnalytics, error) {
	return u.bookingRepo.GetBookingAnalytics(eventID)
}

func (u *bookingUsecaseImpl) validateBooking(booking *booking.Booking) error {
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
