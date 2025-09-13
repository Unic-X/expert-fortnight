package service

import "evently/internal/domain/model"

type NotificationService interface {
	SendBookingConfirmation(user *model.User, booking *model.Booking) error
	SendCancellationNotice(user *model.User, booking *model.Booking) error
	SendWaitlistNotification(user *model.User, event *model.Event) error
}
