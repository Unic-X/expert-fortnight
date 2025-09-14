package di

import (
	"context"

	"evently/internal/config"
	"evently/internal/delivery/http/middleware"
	"evently/internal/domain/model"
	"evently/internal/domain/usecase"
	ucImpl "evently/internal/usecase/impl"
	repoImpl "evently/internal/usecase/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Container struct {
	// Config
	Config *model.Config

	// Database
	Pool *pgxpool.Pool

	// Repositories
	UserRepo    model.UserRepository
	EventRepo   model.EventRepository
	BookingRepo model.BookingRepository

	// Use Cases
	AuthUseCase    usecase.AuthUseCase
	EventUseCase   usecase.EventUsecase
	BookingUseCase usecase.BookingUsecase

	// Middleware
	JWTMiddleware *middleware.JWTConfig

	// Server
	Server *gin.Engine
}

func NewContainer(ctx context.Context) (*Container, error) {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database connection
	pool, err := config.NewPGXPool(ctx, cfg.DB)
	if err != nil {
		return nil, err
	}

	userRepo := repoImpl.NewUserRepository(pool)
	eventRepo := repoImpl.NewEventRepository(pool)
	bookingRepo := repoImpl.NewBookingRepository(pool)

	authUseCase := ucImpl.NewAuthUseCase(userRepo, cfg)
	eventUseCase := ucImpl.NewEventUsecase(eventRepo)
	bookingUseCase := ucImpl.NewBookingUsecase(bookingRepo, eventRepo)

	jwtMiddleware := middleware.NewJWTConfig()

	server := gin.Default()

	return &Container{
		Config:         cfg,
		Pool:           pool,
		UserRepo:       userRepo,
		EventRepo:      eventRepo,
		BookingRepo:    bookingRepo,
		AuthUseCase:    authUseCase,
		EventUseCase:   eventUseCase,
		BookingUseCase: bookingUseCase,
		JWTMiddleware:  jwtMiddleware,
		Server:         server,
	}, nil
}
