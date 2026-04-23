-- +goose Up

CREATE TYPE activity_enum AS ENUM(
	'available', 
	'on_contract',
	'sick_leave',
	'retired',
	'dead'
);

CREATE TYPE role_enum AS ENUM(
	'frontliner',
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

CREATE INDEX ix_ad_rank ON adventurers (current_rank);
CREATE INDEX ix_ad_activity ON adventurers (current_activity);
CREATE INDEX ix_ad_activity ON adventurers (role);
CREATE INDEX ix_party_id ON adventurers (party_id);
CREATE INDEX ix_guild_id ON adventurers (guild_id);

-- +goose Down
DROP INDEX ix_party_id ON adventurers (party_id);
DROP INDEX ix_guild_id ON adventurers (guild_id);
DROP INDEX ix_ad_rank ON adventurers (current_rank);
DROP INDEX ix_ad_activity ON adventurers (current_activity);
DROP INDEX ix_ad_activity ON adventurers (role);
DROP TYPE rank_enum;
DROP TYPE activity_enum;
DROP TYPE role_enum;
DROP TABLE adventurers;
