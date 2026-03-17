package postgres

import (
	"database/sql"

	"github.com/dilvi/camp-booking-rest-api-go/internal/domain"
)

type ChildRepository struct {
	db *sql.DB
}

func NewChildRepository(db *sql.DB) *ChildRepository {
	return &ChildRepository{db: db}
}

func (r *ChildRepository) Create(child domain.Child) (domain.Child, error) {
	query := `
		INSERT INTO children (
			user_id, photo_url, first_name, last_name, birth_date, gender, hobby, allergy
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		child.UserID,
		child.PhotoURL,
		child.FirstName,
		child.LastName,
		child.BirthDate,
		child.Gender,
		child.Hobby,
		child.Allergy,
	).Scan(&child.ID, &child.CreatedAt, &child.UpdatedAt)

	if err != nil {
		return domain.Child{}, err
	}

	return child, nil
}

func (r *ChildRepository) GetAllByUserID(userID int64) ([]domain.Child, error) {
	query := `
		SELECT id, user_id, photo_url, first_name, last_name, birth_date, gender, hobby, allergy, created_at, updated_at
		FROM children
		WHERE user_id = $1
		ORDER BY id
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var children []domain.Child

	for rows.Next() {
		var child domain.Child

		err := rows.Scan(
			&child.ID,
			&child.UserID,
			&child.PhotoURL,
			&child.FirstName,
			&child.LastName,
			&child.BirthDate,
			&child.Gender,
			&child.Hobby,
			&child.Allergy,
			&child.CreatedAt,
			&child.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		children = append(children, child)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return children, nil
}

func (r *ChildRepository) Update(child domain.Child) (domain.Child, error) {
	query := `
		UPDATE children
		SET photo_url = $1,
		    first_name = $2,
		    last_name = $3,
		    birth_date = $4,
		    gender = $5,
		    hobby = $6,
		    allergy = $7,
		    updated_at = NOW()
		WHERE id = $8 AND user_id = $9
		RETURNING updated_at
	`

	err := r.db.QueryRow(
		query,
		child.PhotoURL,
		child.FirstName,
		child.LastName,
		child.BirthDate,
		child.Gender,
		child.Hobby,
		child.Allergy,
		child.ID,
		child.UserID,
	).Scan(&child.UpdatedAt)

	if err != nil {
		return domain.Child{}, err
	}

	return child, nil
}