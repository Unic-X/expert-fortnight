package impl

import (
	"context"
	"fmt"
	"time"

	"evently/internal/domain/events"

	"github.com/google/uuid"
)

type eventUsecaseImpl struct {
	eventRepo events.EventRepository
}

func NewEventUsecase(eventRepo events.EventRepository) events.EventUsecase {
	return &eventUsecaseImpl{
		eventRepo: eventRepo,
	}
}

func (u *eventUsecaseImpl) CreateEvent(ctx context.Context, event *events.Event) error {
	event.ID = uuid.New().String()
	event.CreatedAt = time.Now()
	event.UpdatedAt = time.Now()

	event.AvailableSeats = event.TotalCapacity

	if err := u.validateEvent(event); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return u.eventRepo.Create(event)
}

func (u *eventUsecaseImpl) UpdateEvent(ctx context.Context, event *events.Event) error {
	existingEvent, err := u.eventRepo.GetByID(event.ID)
	if err != nil {
		return fmt.Errorf("event not found: %w", err)
	}

	event.CreatedAt = existingEvent.CreatedAt
	event.CreatedBy = existingEvent.CreatedBy
	event.AvailableSeats = existingEvent.AvailableSeats
	event.CreatedBy = existingEvent.CreatedBy
	event.UpdatedAt = time.Now()

	if err := u.validateEvent(event); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return u.eventRepo.Update(event)
}

func (u *eventUsecaseImpl) DeleteEvent(ctx context.Context, eventID string) error {
	_, err := u.eventRepo.GetByID(eventID)
	if err != nil {
		return fmt.Errorf("event not found: %w", err)
	}

	return u.eventRepo.Delete(eventID)
}

func (u *eventUsecaseImpl) GetEvent(ctx context.Context, eventID string) (*events.Event, error) {
	return u.eventRepo.GetByID(eventID)
}

func (u *eventUsecaseImpl) ListUpcomingEvents(ctx context.Context, limit, offset int) ([]*events.Event, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	return u.eventRepo.ListUpcoming(limit, offset)
}

func (u *eventUsecaseImpl) ListAllEvents(ctx context.Context, limit, offset int) ([]*events.Event, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	return u.eventRepo.ListAll(limit, offset)
}

func (u *eventUsecaseImpl) GetMostPopularEvents(ctx context.Context, limit int) ([]*events.EventAnalytics, error) {
	if limit <= 0 {
		limit = 10
	}

	return u.eventRepo.GetMostPopularEvents(ctx, limit)
}

func (u *eventUsecaseImpl) validateEvent(event *events.Event) error {
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
