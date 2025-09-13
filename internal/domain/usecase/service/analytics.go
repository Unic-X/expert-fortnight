package service

import "evently/internal/domain/model"

type AnalyticsService interface {
	MostPopularEvents(limit int) ([]*model.Event, error)
	CapacityUtilization(eventID string) (float64, error)
	TotalBookings() (int, error)
	DailyBookingStats() (map[string]int, error)
}
