package usecase

import (
	"context"
	"evently/internal/domain/model"
)

type EventUsecase interface {
	CreateEvent(ctx context.Context, event *model.Event) error
	UpdateEvent(ctx context.Context, event *model.Event) error
	DeleteEvent(ctx context.Context, eventID string) error
	GetEvent(ctx context.Context, eventID string) (*model.Event, error)
	ListUpcomingEvents(ctx context.Context, limit, offset int) ([]*model.Event, error)
	ListAllEvents(ctx context.Context, limit, offset int) ([]*model.Event, error)
	GetMostPopularEvents(ctx context.Context, limit int) ([]*model.EventAnalytics, error)
}
