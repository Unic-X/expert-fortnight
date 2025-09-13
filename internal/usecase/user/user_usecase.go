package user

import (
	"evently/internal/domain/model"
	"evently/internal/domain/repository"
	usecase "evently/internal/domain/usecase/user"
)

type userUsecase struct {
	repo repository.EventRepository
}

func NewUserUsecase(repo repository.EventRepository) usecase.UserUsecase {
	return &userUsecase{repo: repo}
}

func (u *userUsecase) GetBookingHistory(userID string) ([]*model.Booking, error) {
	panic("unimplemented")
}

func (u *userUsecase) GetProfile(userID string) (*model.User, error) {
	panic("unimplemented")
}

func (u *userUsecase) Login(email string, password string) (string, error) {
	panic("unimplemented")
}

func (u *userUsecase) Register(name string, email string, password string) (*model.User, error) {
	panic("unimplemented")
}
