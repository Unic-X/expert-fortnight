package usecase

import "evently/internal/domain/model"

type AuthUseCase interface {
	Register(user *model.RegisterRequest) error
	Login(email, password string) (string, error)
}
