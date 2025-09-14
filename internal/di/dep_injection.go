package di

import (
	"context"

	"evently/internal/config"
	"evently/internal/delivery/http/middleware"
	"evently/internal/domain/booking"
	"evently/internal/domain/events"
	"evently/internal/domain/waitlist"

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
	UserRepo         model.UserRepository
	EventRepo        events.EventRepository
	BookingRepo      booking.BookingRepository
	WaitlistRepo     waitlist.WaitlistRepository
	NotificationRepo model.NotificationRepository

	// Use Cases
	AuthUseCase         usecase.AuthUseCase
	EventUseCase        events.EventUsecase
	BookingUseCase      booking.BookingUsecase
	WaitlistUseCase     waitlist.WaitlistUsecase
	NotificationUseCase usecase.NotificationUsecase

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
	waitlistRepo := repoImpl.NewWaitlistRepository(pool)
	notificationRepo := repoImpl.NewNotificationRepository(pool)

	// Initialize use cases
	authUseCase := ucImpl.NewAuthUseCase(userRepo, cfg)
	eventUseCase := ucImpl.NewEventUsecase(eventRepo)
	notificationUseCase := ucImpl.NewNotificationUsecase(notificationRepo, eventRepo)
	waitlistUseCase := ucImpl.NewWaitlistUsecase(waitlistRepo, eventRepo, notificationRepo)
	bookingUseCase := ucImpl.NewBookingUsecase(bookingRepo, eventRepo, waitlistUseCase)

	jwtMiddleware := middleware.NewJWTConfig()

	server := gin.Default()

	return &Container{
		Config:              cfg,
		Pool:                pool,
		UserRepo:            userRepo,
		EventRepo:           eventRepo,
		BookingRepo:         bookingRepo,
		WaitlistRepo:        waitlistRepo,
		NotificationRepo:    notificationRepo,
		AuthUseCase:         authUseCase,
		EventUseCase:        eventUseCase,
		BookingUseCase:      bookingUseCase,
		WaitlistUseCase:     waitlistUseCase,
		NotificationUseCase: notificationUseCase,
		JWTMiddleware:       jwtMiddleware,
		Server:              server,
	}, nil
}
