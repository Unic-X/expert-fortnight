package admin

import "evently/internal/domain/model"

type AdminUsecase interface {
	CreateEvent(event *model.Event) error
	UpdateEvent(event *model.Event) error
	DeleteEvent(eventID string) error
	GetEventAnalytics(eventID string) (map[string]interface{}, error)
	GetOverallAnalytics() (map[string]interface{}, error)
}
