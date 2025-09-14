package admin

import "evently/internal/domain/events"

type AdminUsecase interface {
	CreateEvent(event *events.Event) error
	UpdateEvent(event *events.Event) error
	DeleteEvent(eventID string) error
	GetEventAnalytics(eventID string) (map[string]interface{}, error)
	GetOverallAnalytics() (map[string]interface{}, error)
}
