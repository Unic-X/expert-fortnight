package impl

import (
	"context"
	"fmt"
	"time"

	"evently/internal/domain/model"
	"evently/internal/domain/usecase"

	"github.com/google/uuid"
)

type waitlistUsecaseImpl struct {
	waitlistRepo     model.WaitlistRepository
	eventRepo        model.EventRepository
	notificationRepo model.NotificationRepository
}

func NewWaitlistUsecase(
	waitlistRepo model.WaitlistRepository,
	eventRepo model.EventRepository,
	notificationRepo model.NotificationRepository,
) usecase.WaitlistUsecase {
	return &waitlistUsecaseImpl{
		waitlistRepo:     waitlistRepo,
		eventRepo:        eventRepo,
		notificationRepo: notificationRepo,
	}
}

func (u *waitlistUsecaseImpl) JoinWaitlist(ctx context.Context, userID, eventID string, quantity int) error {
	// Validate input
	if quantity <= 0 {
		return fmt.Errorf("quantity must be positive")
	}
	if quantity > 10 {
		return fmt.Errorf("cannot waitlist for more than 10 tickets at once")
	}

	// Check if event exists
	event, err := u.eventRepo.GetByID(eventID)
	if err != nil {
		return fmt.Errorf("event not found: %w", err)
	}

	// Check if event is in the future
	if event.EventTime.Before(time.Now()) {
		return fmt.Errorf("cannot join waitlist for past events")
	}

	// Check if user is already on waitlist for this event
	existing, err := u.waitlistRepo.GetByUserAndEvent(userID, eventID)
	if err == nil && existing != nil {
		return fmt.Errorf("user is already on waitlist for this event")
	}

	// Create waitlist entry
	waitlist := &model.Waitlist{
		ID:        uuid.New().String(),
		UserID:    userID,
		EventID:   eventID,
		Quantity:  quantity,
		Priority:  0, // Default priority
		Status:    model.WaitlistStatusActive,
		JoinedAt:  time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return u.waitlistRepo.Create(waitlist)
}

func (u *waitlistUsecaseImpl) LeaveWaitlist(ctx context.Context, userID, eventID string) error {
	// Find user's waitlist entry
	waitlist, err := u.waitlistRepo.GetByUserAndEvent(userID, eventID)
	if err != nil {
		return fmt.Errorf("waitlist entry not found: %w", err)
	}

	// Delete the waitlist entry
	return u.waitlistRepo.Delete(waitlist.ID)
}

func (u *waitlistUsecaseImpl) GetUserWaitlist(ctx context.Context, userID string, limit, offset int) ([]*model.Waitlist, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	return u.waitlistRepo.GetByUserID(userID, limit, offset)
}

func (u *waitlistUsecaseImpl) GetEventWaitlist(ctx context.Context, eventID string, limit, offset int) ([]*model.Waitlist, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	return u.waitlistRepo.GetByEventID(eventID, limit, offset)
}

func (u *waitlistUsecaseImpl) ProcessWaitlistNotifications(ctx context.Context, eventID string, availableQuantity int) error {
	// Get next users in queue who can be accommodated
	waitlistEntries, err := u.waitlistRepo.GetNextInQueue(eventID, availableQuantity)
	if err != nil {
		return fmt.Errorf("failed to get waitlist queue: %w", err)
	}

	remainingQuantity := availableQuantity
	for _, entry := range waitlistEntries {
		if remainingQuantity <= 0 {
			break
		}

		if entry.Quantity <= remainingQuantity {
			// Notify this user
			err := u.notifyWaitlistUser(ctx, entry)
			if err != nil {
				// Log error but continue with other users
				fmt.Printf("Failed to notify user %s: %v\n", entry.UserID, err)
				continue
			}

			// Update waitlist status and set expiration
			entry.Status = model.WaitlistStatusNotified
			now := time.Now()
			entry.NotifiedAt = &now
			// Give user 30 minutes to complete booking
			expiresAt := now.Add(30 * time.Minute)
			entry.ExpiresAt = &expiresAt
			entry.UpdatedAt = now

			if err := u.waitlistRepo.Update(entry); err != nil {
				fmt.Printf("Failed to update waitlist entry %s: %v\n", entry.ID, err)
				continue
			}

			remainingQuantity -= entry.Quantity
		}
	}

	return nil
}

func (u *waitlistUsecaseImpl) GetWaitlistPosition(ctx context.Context, userID, eventID string) (int, error) {
	// Get user's waitlist entry
	userEntry, err := u.waitlistRepo.GetByUserAndEvent(userID, eventID)
	if err != nil {
		return 0, fmt.Errorf("user not on waitlist: %w", err)
	}

	// Get all active waitlist entries for this event ordered by priority and join time
	allEntries, err := u.waitlistRepo.GetByEventID(eventID, 1000, 0) // Get a large number
	if err != nil {
		return 0, fmt.Errorf("failed to get waitlist: %w", err)
	}

	// Find user's position
	position := 1
	for _, entry := range allEntries {
		if entry.Status != model.WaitlistStatusActive {
			continue
		}
		if entry.ID == userEntry.ID {
			return position, nil
		}
		position++
	}

	return 0, fmt.Errorf("user position not found")
}

func (u *waitlistUsecaseImpl) ConvertWaitlistToBooking(ctx context.Context, waitlistID string) error {
	// Get waitlist entry
	waitlist, err := u.waitlistRepo.GetByID(waitlistID)
	if err != nil {
		return fmt.Errorf("waitlist entry not found: %w", err)
	}

	// Check if notification is still valid
	if waitlist.Status != model.WaitlistStatusNotified {
		return fmt.Errorf("waitlist entry is not in notified status")
	}

	if waitlist.ExpiresAt != nil && waitlist.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("waitlist notification has expired")
	}

	// Update status to converted
	waitlist.Status = model.WaitlistStatusConverted
	waitlist.UpdatedAt = time.Now()

	return u.waitlistRepo.Update(waitlist)
}

func (u *waitlistUsecaseImpl) CleanupExpiredWaitlist(ctx context.Context) error {
	return u.waitlistRepo.CleanupExpired()
}

func (u *waitlistUsecaseImpl) notifyWaitlistUser(ctx context.Context, waitlist *model.Waitlist) error {
	// Get event details for notification
	event, err := u.eventRepo.GetByID(waitlist.EventID)
	if err != nil {
		return fmt.Errorf("failed to get event details: %w", err)
	}

	// Create notification
	notification := &model.Notification{
		ID:        uuid.New().String(),
		UserID:    waitlist.UserID,
		EventID:   waitlist.EventID,
		Type:      model.NotificationTypeWaitlistSpotAvailable,
		Title:     "Spot Available!",
		Message:   fmt.Sprintf("A spot is now available for %s. You have 30 minutes to complete your booking for %d ticket(s).", event.Name, waitlist.Quantity),
		IsRead:    false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return u.notificationRepo.Create(notification)
}

func (u *waitlistUsecaseImpl) GetWaitlistByID(ctx context.Context, waitlistID string) (*model.Waitlist, error) {
	return u.waitlistRepo.GetByID(waitlistID)
}
