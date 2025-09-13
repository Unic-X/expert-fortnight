package impl

import (
	"context"
	"fmt"
	"time"

	"evently/internal/domain/model"
	"evently/internal/domain/usecase"

	"github.com/google/uuid"
)

type eventUsecaseImpl struct {
	eventRepo model.EventRepository
}

func NewEventUsecase(eventRepo model.EventRepository) usecase.EventUsecase {
	return &eventUsecaseImpl{
		eventRepo: eventRepo,
	}
}

func (u *eventUsecaseImpl) CreateEvent(ctx context.Context, event *model.Event) error {
	// Generate ID and set timestamps
	event.ID = uuid.New().String()
	event.CreatedAt = time.Now()
	event.UpdatedAt = time.Now()

	// Set available seats equal to total capacity initially
	event.AvailableSeats = event.TotalCapacity

	// Validate event data
	if err := u.validateEvent(event); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return u.eventRepo.Create(event)
}

func (u *eventUsecaseImpl) UpdateEvent(ctx context.Context, event *model.Event) error {
	// Check if event exists
	existingEvent, err := u.eventRepo.GetByID(event.ID)
	if err != nil {
		return fmt.Errorf("event not found: %w", err)
	}

	event.CreatedAt = existingEvent.CreatedAt
	event.CreatedBy = existingEvent.CreatedBy
	event.AvailableSeats = existingEvent.AvailableSeats
	event.CreatedBy = existingEvent.CreatedBy
	event.UpdatedAt = time.Now()

	// Validate updated event data
	if err := u.validateEvent(event); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return u.eventRepo.Update(event)
}

func (u *eventUsecaseImpl) DeleteEvent(ctx context.Context, eventID string) error {
	// Check if event exists
	_, err := u.eventRepo.GetByID(eventID)
	if err != nil {
		return fmt.Errorf("event not found: %w", err)
	}

	return u.eventRepo.Delete(eventID)
}

func (u *eventUsecaseImpl) GetEvent(ctx context.Context, eventID string) (*model.Event, error) {
	return u.eventRepo.GetByID(eventID)
}

func (u *eventUsecaseImpl) ListUpcomingEvents(ctx context.Context, limit, offset int) ([]*model.Event, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	return u.eventRepo.ListUpcoming(limit, offset)
}

func (u *eventUsecaseImpl) ListAllEvents(ctx context.Context, limit, offset int) ([]*model.Event, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	return u.eventRepo.ListAll(limit, offset)
}

func (u *eventUsecaseImpl) GetMostPopularEvents(ctx context.Context, limit int) ([]*model.EventAnalytics, error) {
	if limit <= 0 {
		limit = 10
	}

	return u.eventRepo.GetMostPopularEvents(limit)
}

func (u *eventUsecaseImpl) validateEvent(event *model.Event) error {
	if event.Name == "" {
		return fmt.Errorf("event name is required")
	}

	if event.Venue == "" {
		return fmt.Errorf("event venue is required")
	}

	if event.EventTime.Before(time.Now()) {
		return fmt.Errorf("event time must be in the future")
	}

	if event.TotalCapacity <= 0 {
		return fmt.Errorf("event capacity must be positive")
	}

	if event.Price < 0 {
		return fmt.Errorf("event price cannot be negative")
	}

	return nil
}
