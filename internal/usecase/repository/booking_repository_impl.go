package repository

import (
	"context"
	"fmt"

	"evently/internal/domain/booking"

	"github.com/jackc/pgx/v5/pgxpool"
)

type bookingRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewBookingRepository(db *pgxpool.Pool) booking.BookingRepository {
	return &bookingRepositoryImpl{db: db}
}

func (r *bookingRepositoryImpl) Create(newBooking *booking.Booking) error {
	query := `
		INSERT INTO bookings (id, user_id, event_id, quantity, total_amount, 
			status, booking_time, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.Exec(context.Background(), query,
		newBooking.ID, newBooking.UserID, newBooking.EventID, newBooking.Quantity,
		newBooking.TotalAmount, newBooking.Status, newBooking.BookingTime,
		newBooking.CreatedAt, newBooking.UpdatedAt)

	return err
}

func (r *bookingRepositoryImpl) Update(oldBooking *booking.Booking) error {
	query := `
		UPDATE bookings 
		SET quantity = $2, total_amount = $3, status = $4, 
			cancelled_at = $5, updated_at = $6
		WHERE id = $1`

	result, err := r.db.Exec(context.Background(), query,
		oldBooking.ID, oldBooking.Quantity, oldBooking.TotalAmount, oldBooking.Status,
		oldBooking.CancelledAt, oldBooking.UpdatedAt)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("booking not found")
	}

	return nil
}

func (r *bookingRepositoryImpl) GetByID(id string) (*booking.Booking, error) {
	query := `
		SELECT id, user_id, event_id, quantity, total_amount, status, 
			booking_time, cancelled_at, created_at, updated_at
		FROM bookings WHERE id = $1`

	oldBooking := &booking.Booking{}
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&oldBooking.ID, &oldBooking.UserID, &oldBooking.EventID, &oldBooking.Quantity,
		&oldBooking.TotalAmount, &oldBooking.Status, &oldBooking.BookingTime,
		&oldBooking.CancelledAt, &oldBooking.CreatedAt, &oldBooking.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return oldBooking, nil
}

func (r *bookingRepositoryImpl) GetByUserID(userID string, limit, offset int) ([]*booking.Booking, error) {
	query := `
		SELECT id, user_id, event_id, quantity, total_amount, status, 
			booking_time, cancelled_at, created_at, updated_at
		FROM bookings 
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(context.Background(), query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []*booking.Booking
	for rows.Next() {
		booking := &booking.Booking{}
		err := rows.Scan(
			&booking.ID, &booking.UserID, &booking.EventID, &booking.Quantity,
			&booking.TotalAmount, &booking.Status, &booking.BookingTime,
			&booking.CancelledAt, &booking.CreatedAt, &booking.UpdatedAt)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}

	return bookings, rows.Err()
}

func (r *bookingRepositoryImpl) GetByEventID(eventID string, limit, offset int) ([]*booking.Booking, error) {
	query := `
		SELECT id, user_id, event_id, quantity, total_amount, status, 
			booking_time, cancelled_at, created_at, updated_at
		FROM bookings 
		WHERE event_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(context.Background(), query, eventID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []*booking.Booking
	for rows.Next() {
		booking := &booking.Booking{}
		err := rows.Scan(
			&booking.ID, &booking.UserID, &booking.EventID, &booking.Quantity,
			&booking.TotalAmount, &booking.Status, &booking.BookingTime,
			&booking.CancelledAt, &booking.CreatedAt, &booking.UpdatedAt)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}

	return bookings, rows.Err()
}

func (r *bookingRepositoryImpl) CountByEventID(eventID string) (int, error) {
	query := `SELECT COUNT(*) FROM bookings WHERE event_id = $1 AND status = 'confirmed'`

	var count int
	err := r.db.QueryRow(context.Background(), query, eventID).Scan(&count)

	return count, err
}

func (r *bookingRepositoryImpl) GetTotalBookings() (int64, error) {
	query := `SELECT COUNT(*) FROM bookings WHERE status = 'confirmed'`

	var count int64
	err := r.db.QueryRow(context.Background(), query).Scan(&count)

	return count, err
}

func (r *bookingRepositoryImpl) GetBookingAnalytics(eventID string) (*booking.BookingAnalytics, error) {
	query := `
		SELECT 
			event_id,
			COUNT(*) as total_bookings,
			COALESCE(SUM(total_amount), 0) as total_revenue,
			COUNT(CASE WHEN status = 'confirmed' THEN 1 END) as confirmed,
			COUNT(CASE WHEN status = 'cancelled' THEN 1 END) as cancelled,
			COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending
		FROM bookings 
		WHERE event_id = $1
		GROUP BY event_id`

	analytics := &booking.BookingAnalytics{}
	err := r.db.QueryRow(context.Background(), query, eventID).Scan(
		&analytics.EventID, &analytics.TotalBookings, &analytics.TotalRevenue,
		&analytics.Confirmed, &analytics.Cancelled, &analytics.Pending)

	if err != nil {
		return nil, err
	}

	return analytics, nil
}
