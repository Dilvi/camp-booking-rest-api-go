package postgres

import (
	"database/sql"

	"github.com/dilvi/camp-booking-rest-api-go/internal/domain"
)

type BikeRouteRepository struct {
	db *sql.DB
}

func NewBikeRouteRepository(db *sql.DB) *BikeRouteRepository {
	return &BikeRouteRepository{db: db}
}

func (r *BikeRouteRepository) GetAll(campID *int64) ([]domain.BikeRoute, error) {
	query := `
		SELECT id, camp_id, title, location, map_url, route_points, distance_km,
		       duration_minutes, difficulty, route_type, elevation_gain_m, description,
		       recommendations_context, created_at, updated_at
		FROM bike_routes
	`

	args := []any{}
	if campID != nil {
		query += " WHERE camp_id = $1"
		args = append(args, *campID)
	}
	query += " ORDER BY id"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var routes []domain.BikeRoute
	for rows.Next() {
		var route domain.BikeRoute
		if err := rows.Scan(
			&route.ID,
			&route.CampID,
			&route.Title,
			&route.Location,
			&route.MapURL,
			&route.RoutePoints,
			&route.DistanceKM,
			&route.DurationMinutes,
			&route.Difficulty,
			&route.RouteType,
			&route.ElevationGainM,
			&route.Description,
			&route.RecommendationsContext,
			&route.CreatedAt,
			&route.UpdatedAt,
		); err != nil {
			return nil, err
		}
		routes = append(routes, route)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return routes, nil
}

func (r *BikeRouteRepository) GetByID(id int64) (domain.BikeRoute, error) {
	query := `
		SELECT id, camp_id, title, location, map_url, route_points, distance_km,
		       duration_minutes, difficulty, route_type, elevation_gain_m, description,
		       recommendations_context, created_at, updated_at
		FROM bike_routes
		WHERE id = $1
	`

	var route domain.BikeRoute
	err := r.db.QueryRow(query, id).Scan(
		&route.ID,
		&route.CampID,
		&route.Title,
		&route.Location,
		&route.MapURL,
		&route.RoutePoints,
		&route.DistanceKM,
		&route.DurationMinutes,
		&route.Difficulty,
		&route.RouteType,
		&route.ElevationGainM,
		&route.Description,
		&route.RecommendationsContext,
		&route.CreatedAt,
		&route.UpdatedAt,
	)
	if err != nil {
		return domain.BikeRoute{}, err
	}

	return route, nil
}
