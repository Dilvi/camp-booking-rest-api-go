CREATE TABLE camps (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    location TEXT NOT NULL,
    price_per_day INTEGER NOT NULL,
    booked_count INTEGER NOT NULL DEFAULT 0,
    description TEXT NOT NULL,
    shift_duration_days INTEGER NOT NULL,
    age_min INTEGER NOT NULL,
    age_max INTEGER NOT NULL,
    camp_type TEXT NOT NULL,
    food_type TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);