package repository

import (
	"evently/internal/domain/model"
	"evently/internal/domain/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type eventRepository struct {
	pool *pgxpool.Pool
}

func NewEventRepository(pool *pgxpool.Pool) repository.EventRepository {
	return &eventRepository{pool: pool}
}

// Admin Only
func (e *eventRepository) Create(event *model.Event) error {
	panic("unimplemented")
}

func (e *eventRepository) DecrementCapacity(eventID string) error {
	panic("unimplemented")
}

func (e *eventRepository) IncrementCapacity(eventID string) error {
	panic("unimplemented")
}

func (e *eventRepository) Delete(id string) error {
	panic("unimplemented")
}

func (e *eventRepository) Update(event *model.Event) error {
	panic("unimplemented")
}

func (e *eventRepository) GetEventAnalytics(eventID string) (map[string]interface{}, error) {
	panic("unimplemented")
}

func (e *eventRepository) GetOverallAnalytics() (map[string]interface{}, error) {
	panic("unimplemented")
}

// User
func (e *eventRepository) GetByID(id string) (*model.Event, error) {
	panic("unimplemented")
}

func (e *eventRepository) ListUpcoming() ([]*model.Event, error) {
	panic("unimplemented")
}
