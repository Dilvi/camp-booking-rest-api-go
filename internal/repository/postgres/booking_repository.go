package postgres

import (
	"database/sql"

	"github.com/dilvi/camp-booking-rest-api-go/internal/domain"
)

type BookingRepository struct {
	db *sql.DB
}

func NewBookingRepository(db *sql.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) Create(b domain.Booking) (domain.Booking, error) {
	query := `
		INSERT INTO bookings (user_id, child_id, camp_id)
		VALUES ($1, $2, $3)
		RETURNING id, status, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		b.UserID,
		b.ChildID,
		b.CampID,
	).Scan(&b.ID, &b.Status, &b.CreatedAt, &b.UpdatedAt)

	if err != nil {
		return domain.Booking{}, err
	}

	return b, nil
}

func (r *BookingRepository) GetAllByUserID(userID int64) ([]domain.Booking, error) {
	query := `
		SELECT id, user_id, child_id, camp_id, status, created_at, updated_at
		FROM bookings
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []domain.Booking

	for rows.Next() {
		var b domain.Booking

		err := rows.Scan(
			&b.ID,
			&b.UserID,
			&b.ChildID,
			&b.CampID,
			&b.Status,
			&b.CreatedAt,
			&b.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		bookings = append(bookings, b)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return bookings, nil
}

func (r *BookingRepository) GetByIDForUser(id, userID int64) (domain.Booking, error) {
	query := `
		SELECT id, user_id, child_id, camp_id, status, created_at, updated_at
		FROM bookings
		WHERE id = $1 AND user_id = $2
	`

	var b domain.Booking
	err := r.db.QueryRow(query, id, userID).Scan(
		&b.ID,
		&b.UserID,
		&b.ChildID,
		&b.CampID,
		&b.Status,
		&b.CreatedAt,
		&b.UpdatedAt,
	)
	if err != nil {
		return domain.Booking{}, err
	}

	return b, nil
}

func (r *BookingRepository) ExistsActiveByUserAndCamp(userID, campID int64) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM bookings
			WHERE user_id = $1
			  AND camp_id = $2
			  AND status IN ('pending', 'confirmed')
		)
	`

	var exists bool
	if err := r.db.QueryRow(query, userID, campID).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}
