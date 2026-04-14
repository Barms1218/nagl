-- +goose Up
CREATE TYPE rank_enum AS ENUM('junior', 'senior');

CREATE TYPE activity_enum AS ENUM(
	'available', 
	'on_quest',
	'sick_leave',
	'retired',
	'dead'
);

CREATE TYPE role_enum AS ENUM(
	'frontline',
	'spellcaster',
	'healer',
	'generalist'
);

CREATE TABLE adventurers (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	guild_id UUID REFERENCES guilds(id),
	joined_at TIMESTAMP  NOT NULL DEFAULT CURRENT_TIMESTAMP,
	recruitment_cost INTEGER NOT NULL DEFAULT 0,
	current_rank rank_enum NOT NULL DEFAULT 'junior',
	current_activity activity_enum NOT NULL DEFAULT 'available',
	name TEXT,
	description TEXT NOT NULL,
	role role_enum NOT NULL DEFAULT 'generalist',
	upkeep_cost INTEGER NOT NULL DEFAULT 0
);

-- +goose Down
DROP TYPE rank_enum;
DROP TYPE activity_enum;
DROP TYPE role_enum;
DROP TABLE adventurers;
