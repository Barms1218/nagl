-- +goose Up
CREATE TYPE rank_enum AS ENUM('junior', 'senior')
CREATE TYPE activity_enum AS ENUM(
	'available', 
	'on_quest',
	'traveling',
	'sick',
	'retired',
	'dead'
)
CREATE TABLE adventurers (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	joined_at TIMESTAMP  NOT NULL DEFAULT CURRENT_TIMESTAMP,
	rank rank_emum NOT NULL DEFAULT 'junior',
	activity activity_enum NOT NULL DEFAULT 'available',
	name TEXT,
	class TEXT,
	stats JSONB NOT NULL DEFAULT '{}'::jsonb,
)

-- +goose Down
DROP TYPE rank_enum;
DROP TYPE status_enum;
DROP TABLE adventurers;
