package user

import "evently/internal/domain/model"

type UserUsecase interface {
	Register(name, email, password string) (*model.User, error)
	Login(email, password string) (string, error) // returns JWT
	GetProfile(userID string) (*model.User, error)
	GetBookingHistory(userID string) ([]*model.Booking, error)
}
