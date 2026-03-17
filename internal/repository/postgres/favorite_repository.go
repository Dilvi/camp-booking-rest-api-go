package postgres

import (
	"database/sql"

	"github.com/dilvi/camp-booking-rest-api-go/internal/domain"
)

type FavoriteRepository struct {
	db *sql.DB
}

func NewFavoriteRepository(db *sql.DB) *FavoriteRepository {
	return &FavoriteRepository{db: db}
}

func (r *FavoriteRepository) Add(userID, campID int64) error {
	query := `
		INSERT INTO favorites (user_id, camp_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`
	_, err := r.db.Exec(query, userID, campID)
	return err
}

func (r *FavoriteRepository) Remove(userID, campID int64) error {
	query := `
		DELETE FROM favorites
		WHERE user_id = $1 AND camp_id = $2
	`
	_, err := r.db.Exec(query, userID, campID)
	return err
}

func (r *FavoriteRepository) GetAll(userID int64) ([]domain.Camp, error) {
	query := `
		SELECT c.id, c.title, c.location, c.price_per_day, c.booked_count,
		       c.description, c.shift_duration_days, c.age_min, c.age_max,
		       c.camp_type, c.food_type, c.created_at, c.updated_at
		FROM favorites f
		JOIN camps c ON c.id = f.camp_id
		WHERE f.user_id = $1
		ORDER BY f.created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var camps []domain.Camp

	for rows.Next() {
		var camp domain.Camp
		err := rows.Scan(
			&camp.ID,
			&camp.Title,
			&camp.Location,
			&camp.PricePerDay,
			&camp.BookedCount,
			&camp.Description,
			&camp.ShiftDurationDays,
			&camp.AgeMin,
			&camp.AgeMax,
			&camp.CampType,
			&camp.FoodType,
			&camp.CreatedAt,
			&camp.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		camps = append(camps, camp)
	}

	return camps, nil
}