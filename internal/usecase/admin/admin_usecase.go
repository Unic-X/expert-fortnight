package admin

import (
	"evently/internal/domain/model"
	"evently/internal/domain/repository"
	usecase "evently/internal/domain/usecase/admin"
)

type adminUsecase struct {
	repo repository.EventRepository
}

func NewAdminUsecase(repo repository.EventRepository) usecase.AdminUsecase {
	return &adminUsecase{repo: repo}
}

func (a *adminUsecase) CreateEvent(event *model.Event) error {
	panic("unimplemented")
}

func (a *adminUsecase) DeleteEvent(eventID string) error {
	panic("unimplemented")
}

func (a *adminUsecase) GetEventAnalytics(eventID string) (map[string]interface{}, error) {
	panic("unimplemented")
}

func (a *adminUsecase) GetOverallAnalytics() (map[string]interface{}, error) {
	panic("unimplemented")
}

func (a *adminUsecase) UpdateEvent(event *model.Event) error {
	panic("unimplemented")
}
