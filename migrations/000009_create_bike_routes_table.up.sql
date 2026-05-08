CREATE TABLE bike_routes (
    id BIGSERIAL PRIMARY KEY,
    camp_id BIGINT NOT NULL REFERENCES camps(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    location TEXT NOT NULL,
    map_url TEXT NOT NULL DEFAULT '',
    route_points TEXT NOT NULL DEFAULT '',
    distance_km NUMERIC(6, 2) NOT NULL,
    duration_minutes INTEGER NOT NULL,
    difficulty TEXT NOT NULL,
    route_type TEXT NOT NULL,
    elevation_gain_m INTEGER NOT NULL DEFAULT 0,
    description TEXT NOT NULL,
    recommendations_context TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
