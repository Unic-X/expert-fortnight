package di

import (
	"context"

	"evently/internal/config"
	"evently/internal/delivery/http/middleware"
	"evently/internal/domain/model"
	"evently/internal/domain/usecase"
	repoImpl "evently/internal/repository/impl"
	ucImpl "evently/internal/usecase/impl"

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

	// Initialize repositories
	userRepo := repoImpl.NewUserRepository(pool)
	eventRepo := repoImpl.NewEventRepository(pool)
	bookingRepo := repoImpl.NewBookingRepository(pool)

	// Initialize use cases
	authUseCase := ucImpl.NewAuthUseCase(userRepo, cfg)
	eventUseCase := ucImpl.NewEventUsecase(eventRepo)
	bookingUseCase := ucImpl.NewBookingUsecase(bookingRepo, eventRepo)

	// Initialize middleware
	jwtMiddleware := middleware.NewJWTConfig()

	// Initialize server
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
