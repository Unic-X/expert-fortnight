package repository

import (
	"context"
	"fmt"

	"evently/internal/domain/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type notificationRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewNotificationRepository(db *pgxpool.Pool) model.NotificationRepository {
	return &notificationRepositoryImpl{db: db}
}

func (r *notificationRepositoryImpl) Create(notification *model.Notification) error {
	query := `
		INSERT INTO notifications (id, user_id, event_id, type, title, message, 
			is_read, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.Exec(context.Background(), query,
		notification.ID, notification.UserID, notification.EventID, notification.Type,
		notification.Title, notification.Message, notification.IsRead,
		notification.CreatedAt, notification.UpdatedAt)

	return err
}

func (r *notificationRepositoryImpl) GetByUserID(userID string, limit, offset int) ([]*model.Notification, error) {
	query := `
		SELECT id, user_id, event_id, type, title, message, is_read, created_at, updated_at
		FROM notifications 
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(context.Background(), query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*model.Notification
	for rows.Next() {
		notification := &model.Notification{}
		err := rows.Scan(
			&notification.ID, &notification.UserID, &notification.EventID, &notification.Type,
			&notification.Title, &notification.Message, &notification.IsRead,
			&notification.CreatedAt, &notification.UpdatedAt)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}

	return notifications, rows.Err()
}

func (r *notificationRepositoryImpl) MarkAsRead(id string) error {
	query := `UPDATE notifications SET is_read = true WHERE id = $1`

	result, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("notification not found")
	}

	return nil
}

func (r *notificationRepositoryImpl) MarkAllAsRead(userID string) error {
	query := `UPDATE notifications SET is_read = true WHERE user_id = $1`

	_, err := r.db.Exec(context.Background(), query, userID)
	return err
}

func (r *notificationRepositoryImpl) GetUnreadCount(userID string) (int, error) {
	query := `SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND is_read = false`

	var count int
	err := r.db.QueryRow(context.Background(), query, userID).Scan(&count)

	return count, err
}
