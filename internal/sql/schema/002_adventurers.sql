-- +goose Up
CREATE TYPE rank_enum AS ENUM('junior', 'senior');

CREATE TYPE activity_enum AS ENUM(
	'available', 
	'on_quest',
	'traveling',
	'sick',
	'retired',
	'dead'
);

CREATE TYPE role_enum AS ENUM(
	'frontline',
	'spellcaster',
	'generalist',
	'healer'
);

CREATE TABLE adventurers (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	guild_id UUID REFERENCES guilds(id),
	joined_at TIMESTAMP  NOT NULL DEFAULT CURRENT_TIMESTAMP,
	recruitment_cost INTEGER NOT NULL DEFAULT 0,
	current_rank rank_emum NOT NULL DEFAULT 'junior',
	current_activity activity_enum NOT NULL DEFAULT 'available',
	name TEXT,
	role role_enum NOT NULL DEFAULT 'generatlist'
);

-- +goose Down
DROP TYPE rank_enum;
DROP TYPE status_enum;
DROP TYPE role_enum;
DROP TABLE adventurers;
