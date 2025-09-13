package impl

import (
	"fmt"
	"log"
	"time"

	"evently/internal/delivery/http/middleware"
	"evently/internal/domain/model"
	"evently/internal/domain/usecase"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type authUsecaseImpl struct {
	userRepo model.UserRepository
	config   *model.Config
}

func NewAuthUseCase(userRepo model.UserRepository, config *model.Config) usecase.AuthUseCase {
	return &authUsecaseImpl{
		userRepo: userRepo,
		config:   config,
	}
}

func (u *authUsecaseImpl) Register(req *model.RegisterRequest) error {
	// Check if user already exists

	newUser := model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	}

	existingUser, _ := u.userRepo.GetByEmail(newUser.Email)
	if existingUser != nil {
		return fmt.Errorf("user with email %s already exists", newUser.Email)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Set user fields
	newUser.ID = uuid.New().String()
	newUser.Password = string(hashedPassword)
	newUser.CreatedAt = time.Now()

	// Set default role if not specified
	if newUser.Role == "" {
		newUser.Role = "user"
	}

	// Validate user data
	if err := u.validateUser(&newUser); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return u.userRepo.Create(&newUser)
}

func (u *authUsecaseImpl) Login(email, password string) (string, error) {
	// Get user by email
	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		return "", fmt.Errorf("email doesn't exist in database please register")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Printf("Password comparison failed: %v", err)
		return "", fmt.Errorf("invalid password")
	}

	// Generate JWT token
	jwtConfig := middleware.NewJWTConfig()
	token, err := jwtConfig.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

func (u *authUsecaseImpl) validateUser(user *model.User) error {
	if user.Name == "" {
		return fmt.Errorf("name is required")
	}

	if user.Email == "" {
		return fmt.Errorf("email is required")
	}

	if user.Role != "user" && user.Role != "admin" {
		return fmt.Errorf("invalid role: must be 'user' or 'admin'")
	}

	return nil
}
