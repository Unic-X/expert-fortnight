package repository

import "evently/internal/domain/model"

type EventRepository interface {
	Create(event *model.Event) error
	Update(event *model.Event) error
	Delete(id string) error
	GetByID(id string) (*model.Event, error)
	ListUpcoming() ([]*model.Event, error)
	DecrementCapacity(eventID string) error
	IncrementCapacity(eventID string) error
	GetEventAnalytics(eventID string) (map[string]interface{}, error)
	GetOverallAnalytics() (map[string]interface{}, error)
}
