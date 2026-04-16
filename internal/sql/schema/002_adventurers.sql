-- +goose Up

CREATE TYPE activity_enum AS ENUM(
	'available', 
	'on_contract',
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
	-- Primary Key
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

	-- Foreign Keys
	guild_id UUID REFERENCES guilds(id) ON DELETE SET NULL,
	party_id UUID REFERENCES parties(id) ON DELETE SET NULL,

	-- Audit Info
	joined_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	
	-- Personal Details
	current_rank INTEGER NOT NULL DEFAULT 1,
	current_activity activity_enum NOT NULL DEFAULT 'available',
	name TEXT,
	description TEXT NOT NULL,
	role role_enum NOT NULL DEFAULT 'generalist',

	-- Financial Data
	upkeep_cost INTEGER NOT NULL DEFAULT 0,
	recruitment_cost INTEGER NOT NULL DEFAULT 0
);

-- +goose Down
DROP TYPE rank_enum;
DROP TYPE activity_enum;
DROP TYPE role_enum;
DROP TABLE adventurers;
