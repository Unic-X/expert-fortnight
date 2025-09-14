package repository

import (
	"context"
	"fmt"
	"time"

	"evently/internal/domain/waitlist"

	"github.com/jackc/pgx/v5/pgxpool"
)

type waitlistRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewWaitlistRepository(db *pgxpool.Pool) waitlist.WaitlistRepository {
	return &waitlistRepositoryImpl{db: db}
}

func (r *waitlistRepositoryImpl) Create(waitlist *waitlist.Waitlist) error {
	query := `
		INSERT INTO waitlist (id, user_id, event_id, quantity, priority, status, 
			joined_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.Exec(context.Background(), query,
		waitlist.ID, waitlist.UserID, waitlist.EventID, waitlist.Quantity,
		waitlist.Priority, waitlist.Status, waitlist.JoinedAt,
		waitlist.CreatedAt, waitlist.UpdatedAt)

	return err
}

func (r *waitlistRepositoryImpl) Update(waitlist *waitlist.Waitlist) error {
	query := `
		UPDATE waitlist 
		SET quantity = $2, priority = $3, status = $4, notified_at = $5, 
			expires_at = $6, updated_at = $7
		WHERE id = $1`

	result, err := r.db.Exec(context.Background(), query,
		waitlist.ID, waitlist.Quantity, waitlist.Priority, waitlist.Status,
		waitlist.NotifiedAt, waitlist.ExpiresAt, waitlist.UpdatedAt)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("waitlist entry not found")
	}

	return nil
}

func (r *waitlistRepositoryImpl) Delete(id string) error {
	query := `DELETE FROM waitlist WHERE id = $1`

	result, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("waitlist entry not found")
	}

	return nil
}

func (r *waitlistRepositoryImpl) GetByID(id string) (*waitlist.Waitlist, error) {
	query := `
		SELECT id, user_id, event_id, quantity, priority, status, 
			joined_at, notified_at, expires_at, created_at, updated_at
		FROM waitlist WHERE id = $1`

	waitlist := &waitlist.Waitlist{}
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&waitlist.ID, &waitlist.UserID, &waitlist.EventID, &waitlist.Quantity,
		&waitlist.Priority, &waitlist.Status, &waitlist.JoinedAt,
		&waitlist.NotifiedAt, &waitlist.ExpiresAt, &waitlist.CreatedAt, &waitlist.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return waitlist, nil
}

func (r *waitlistRepositoryImpl) GetByUserAndEvent(userID, eventID string) (*waitlist.Waitlist, error) {
	query := `
		SELECT id, user_id, event_id, quantity, priority, status, 
			joined_at, notified_at, expires_at, created_at, updated_at
		FROM waitlist WHERE user_id = $1 AND event_id = $2`

	waitlist := &waitlist.Waitlist{}
	err := r.db.QueryRow(context.Background(), query, userID, eventID).Scan(
		&waitlist.ID, &waitlist.UserID, &waitlist.EventID, &waitlist.Quantity,
		&waitlist.Priority, &waitlist.Status, &waitlist.JoinedAt,
		&waitlist.NotifiedAt, &waitlist.ExpiresAt, &waitlist.CreatedAt, &waitlist.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return waitlist, nil
}

func (r *waitlistRepositoryImpl) GetByEventID(eventID string, limit, offset int) ([]*waitlist.Waitlist, error) {
	query := `
		SELECT id, user_id, event_id, quantity, priority, status, 
			joined_at, notified_at, expires_at, created_at, updated_at
		FROM waitlist 
		WHERE event_id = $1
		ORDER BY priority DESC, joined_at ASC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(context.Background(), query, eventID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var waitlists []*waitlist.Waitlist
	for rows.Next() {
		waitlist := &waitlist.Waitlist{}
		err := rows.Scan(
			&waitlist.ID, &waitlist.UserID, &waitlist.EventID, &waitlist.Quantity,
			&waitlist.Priority, &waitlist.Status, &waitlist.JoinedAt,
			&waitlist.NotifiedAt, &waitlist.ExpiresAt, &waitlist.CreatedAt, &waitlist.UpdatedAt)
		if err != nil {
			return nil, err
		}
		waitlists = append(waitlists, waitlist)
	}

	return waitlists, rows.Err()
}

func (r *waitlistRepositoryImpl) GetByUserID(userID string, limit, offset int) ([]*waitlist.Waitlist, error) {
	query := `
		SELECT id, user_id, event_id, quantity, priority, status, 
			joined_at, notified_at, expires_at, created_at, updated_at
		FROM waitlist 
		WHERE user_id = $1
		ORDER BY joined_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(context.Background(), query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var waitlists []*waitlist.Waitlist
	for rows.Next() {
		waitlist := &waitlist.Waitlist{}
		err := rows.Scan(
			&waitlist.ID, &waitlist.UserID, &waitlist.EventID, &waitlist.Quantity,
			&waitlist.Priority, &waitlist.Status, &waitlist.JoinedAt,
			&waitlist.NotifiedAt, &waitlist.ExpiresAt, &waitlist.CreatedAt, &waitlist.UpdatedAt)
		if err != nil {
			return nil, err
		}
		waitlists = append(waitlists, waitlist)
	}

	return waitlists, rows.Err()
}

func (r *waitlistRepositoryImpl) GetNextInQueue(eventID string, quantity int) ([]*waitlist.Waitlist, error) {
	query := `
		SELECT id, user_id, event_id, quantity, priority, status, 
			joined_at, notified_at, expires_at, created_at, updated_at
		FROM waitlist 
		WHERE event_id = $1 AND status = 'active' AND quantity <= $2
		ORDER BY priority DESC, joined_at ASC`

	rows, err := r.db.Query(context.Background(), query, eventID, quantity)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var waitlists []*waitlist.Waitlist
	for rows.Next() {
		waitlist := &waitlist.Waitlist{}
		err := rows.Scan(
			&waitlist.ID, &waitlist.UserID, &waitlist.EventID, &waitlist.Quantity,
			&waitlist.Priority, &waitlist.Status, &waitlist.JoinedAt,
			&waitlist.NotifiedAt, &waitlist.ExpiresAt, &waitlist.CreatedAt, &waitlist.UpdatedAt)
		if err != nil {
			return nil, err
		}
		waitlists = append(waitlists, waitlist)
	}

	return waitlists, rows.Err()
}

func (r *waitlistRepositoryImpl) CountByEventID(eventID string) (int, error) {
	query := `SELECT COUNT(*) FROM waitlist WHERE event_id = $1 AND status = 'active'`

	var count int
	err := r.db.QueryRow(context.Background(), query, eventID).Scan(&count)

	return count, err
}

func (r *waitlistRepositoryImpl) UpdateStatus(id string, status waitlist.WaitlistStatus) error {
	query := `UPDATE waitlist SET status = $2, updated_at = $3 WHERE id = $1`

	result, err := r.db.Exec(context.Background(), query, id, status, time.Now())
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("waitlist entry not found")
	}

	return nil
}

func (r *waitlistRepositoryImpl) CleanupExpired() error {
	query := `
		UPDATE waitlist 
		SET status = 'expired', updated_at = $1 
		WHERE status = 'notified' AND expires_at < $1`

	_, err := r.db.Exec(context.Background(), query, time.Now())
	return err
}
