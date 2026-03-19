package postgres

import (
	"database/sql"

	"github.com/dilvi/camp-booking-rest-api-go/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user domain.User) (domain.User, error) {
	query := `
		INSERT INTO users (first_name, last_name, phone, email, password_hash, avatar_url, role)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		user.FirstName,
		user.LastName,
		user.Phone,
		user.Email,
		user.PasswordHash,
		user.AvatarURL,
		user.Role,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (r *UserRepository) GetByEmail(email string) (domain.User, error) {
	query := `
		SELECT id, first_name, last_name, phone, email, password_hash, avatar_url, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user domain.User

	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.Email,
		&user.PasswordHash,
		&user.AvatarURL,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (r *UserRepository) GetByID(id int64) (domain.User, error) {
	query := `
		SELECT id, first_name, last_name, phone, email, password_hash, avatar_url, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user domain.User

	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.Email,
		&user.PasswordHash,
		&user.AvatarURL,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (r *UserRepository) UpdateProfile(user domain.User) (domain.User, error) {
	query := `
		UPDATE users
		SET first_name = $1,
		    last_name = $2,
		    phone = $3,
		    email = $4,
		    avatar_url = $5,
		    updated_at = NOW()
		WHERE id = $6
		RETURNING updated_at
	`

	err := r.db.QueryRow(
		query,
		user.FirstName,
		user.LastName,
		user.Phone,
		user.Email,
		user.AvatarURL,
		user.ID,
	).Scan(&user.UpdatedAt)

	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (r *UserRepository) UpdatePassword(userID int64, passwordHash string) error {
	query := `
		UPDATE users
		SET password_hash = $1,
		    updated_at = NOW()
		WHERE id = $2
	`

	_, err := r.db.Exec(query, passwordHash, userID)
	return err
}