package postgres

import (
	"database/sql"

	"github.com/dilvi/camp-booking-rest-api-go/internal/domain"
)

type CampRepository struct {
	db *sql.DB
}

func NewCampRepository(db *sql.DB) *CampRepository {
	return &CampRepository{db: db}
}

func (r *CampRepository) GetAll() ([]domain.Camp, error) {
	query := `
		SELECT id, title, location, image_url, price_per_day, booked_count, description,
			shift_duration_days, age_min, age_max, camp_type, food_type,
			created_at, updated_at
		FROM camps
		ORDER BY id
	`

	rows, err := r.db.Query(query)
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
			&camp.ImageURL,
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

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return camps, nil
}

func (r *CampRepository) GetByID(id int64) (domain.Camp, error) {
	query := `
		SELECT id, title, location, image_url, price_per_day, booked_count, description,
			shift_duration_days, age_min, age_max, camp_type, food_type,
			created_at, updated_at
		FROM camps
		WHERE id = $1
	`

	var camp domain.Camp

	err := r.db.QueryRow(query, id).Scan(
		&camp.ID,
		&camp.Title,
		&camp.Location,
		&camp.ImageURL,
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
		return domain.Camp{}, err
	}

	return camp, nil
}