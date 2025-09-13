package service

import "evently/internal/domain/model"

type AuthService interface {
	HashPassword(password string) (string, error)
	VerifyPassword(hashed, plain string) bool
	GenerateToken(user *model.User) (string, error)
	ValidateToken(token string) (*model.User, error)
}
