package repository

import (
	"context"
	"fmt"
	"time"

	"evently/internal/domain/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type eventRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewEventRepository(db *pgxpool.Pool) model.EventRepository {
	return &eventRepositoryImpl{db: db}
}

func (r *eventRepositoryImpl) Create(event *model.Event) error {
	query := `
		INSERT INTO events (id, name, description, venue, event_time, total_capacity, 
			available_seats, price, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := r.db.Exec(context.Background(), query,
		event.ID, event.Name, event.Description, event.Venue, event.EventTime,
		event.TotalCapacity, event.AvailableSeats, event.Price, event.CreatedBy,
		event.CreatedAt, event.UpdatedAt)

	return err
}

func (r *eventRepositoryImpl) Update(event *model.Event) error {
	query := `
		UPDATE events 
		SET name = $2, description = $3, venue = $4, event_time = $5, 
			total_capacity = $6, available_seats = $7, price = $8, updated_at = $9
		WHERE id = $1`

	result, err := r.db.Exec(context.Background(), query,
		event.ID, event.Name, event.Description, event.Venue, event.EventTime,
		event.TotalCapacity, event.AvailableSeats, event.Price, event.UpdatedAt)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("event not found")
	}

	return nil
}

func (r *eventRepositoryImpl) Delete(id string) error {
	query := `DELETE FROM events WHERE id = $1`

	result, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("event not found")
	}

	return nil
}

func (r *eventRepositoryImpl) GetByID(id string) (*model.Event, error) {
	query := `
		SELECT id, name, description, venue, event_time, total_capacity, 
			available_seats, price, created_by, created_at, updated_at
		FROM events WHERE id = $1`

	event := &model.Event{}
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&event.ID, &event.Name, &event.Description, &event.Venue, &event.EventTime,
		&event.TotalCapacity, &event.AvailableSeats, &event.Price, &event.CreatedBy,
		&event.CreatedAt, &event.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return event, nil
}

func (r *eventRepositoryImpl) ListUpcoming(limit, offset int) ([]*model.Event, error) {
	query := `
		SELECT id, name, description, venue, event_time, total_capacity, 
			available_seats, price, created_by, created_at, updated_at
		FROM events 
		WHERE event_time > NOW()
		ORDER BY event_time ASC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(context.Background(), query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*model.Event
	for rows.Next() {
		event := &model.Event{}
		err := rows.Scan(
			&event.ID, &event.Name, &event.Description, &event.Venue, &event.EventTime,
			&event.TotalCapacity, &event.AvailableSeats, &event.Price, &event.CreatedBy,
			&event.CreatedAt, &event.UpdatedAt)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, rows.Err()
}

func (r *eventRepositoryImpl) ListAll(limit, offset int) ([]*model.Event, error) {
	query := `
		SELECT id, name, description, venue, event_time, total_capacity, 
			available_seats, price, created_by, created_at, updated_at
		FROM events 
		ORDER BY event_time DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(context.Background(), query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*model.Event
	for rows.Next() {
		event := &model.Event{}
		err := rows.Scan(
			&event.ID, &event.Name, &event.Description, &event.Venue, &event.EventTime,
			&event.TotalCapacity, &event.AvailableSeats, &event.Price, &event.CreatedBy,
			&event.CreatedAt, &event.UpdatedAt)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, rows.Err()
}

func (r *eventRepositoryImpl) UpdateAvailableSeats(eventID string, quantity int) error {
	query := `
		UPDATE events 
		SET available_seats = available_seats + $2, updated_at = $3
		WHERE id = $1 AND available_seats + $2 >= 0`

	result, err := r.db.Exec(context.Background(), query, eventID, quantity, time.Now())
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("insufficient seats or event not found")
	}

	return nil
}

func (r *eventRepositoryImpl) GetMostPopularEvents(limit int) ([]*model.EventAnalytics, error) {
	query := `
		SELECT 
			e.id as event_id,
			e.name as event_name,
			COUNT(b.id) as total_bookings,
			COALESCE(SUM(b.total_amount), 0) as total_revenue,
			COALESCE(SUM(b.quantity), 0) as capacity_used,
			e.total_capacity as capacity_total
		FROM events e
		LEFT JOIN bookings b ON e.id = b.event_id AND b.status = 'confirmed'
		GROUP BY e.id, e.name, e.total_capacity
		ORDER BY total_bookings DESC, total_revenue DESC
		LIMIT $1`

	rows, err := r.db.Query(context.Background(), query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var analytics []*model.EventAnalytics
	for rows.Next() {
		analytic := &model.EventAnalytics{}
		err := rows.Scan(
			&analytic.EventID, &analytic.EventName, &analytic.TotalBookings,
			&analytic.TotalRevenue, &analytic.CapacityUsed, &analytic.CapacityTotal)
		if err != nil {
			return nil, err
		}

		// Calculate utilization rate
		if analytic.CapacityTotal > 0 {
			analytic.UtilizationRate = float64(analytic.CapacityUsed) / float64(analytic.CapacityTotal) * 100
		}

		analytics = append(analytics, analytic)
	}

	return analytics, rows.Err()
}
