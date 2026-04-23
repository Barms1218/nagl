-- +goose Up

CREATE TYPE party_status_enum AS ENUM(
	'available',
	'training',
	'fighting', 
	'camping',
	'eating',
	'traveling',
	'dead'
);

CREATE TABLE parties (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	contract_id UUID REFERENCES contracts(id) ON DELETE SET NULL,
	guild_id UUID NOT NULL REFERENCES guilds(id) ON DELETE CASCADE,
	name TEXT UNIQUE NOT NULL,
	party_rank INTEGER NOT NULL DEFAULT 1,
	maximum_party_size INTEGER NOT NULL DEFAULT 3,
	party_status party_status_enum NOT NULL DEFAULT 'available',
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UNIQUE(guild_id, contract_id)
);

CREATE INDEX ix_party_status ON parties (party_status);
CREATE INDEX ix_party_guild ON parties (guild_id);
CREATE INDEX ix_party_name ON parties (name);
CREATE INDEX ix_party_rank ON parties (party_rank);
CREATE INDEX ix_party_size ON parties (maximum_party_size);

-- +goose Down
DROP INDEX ix_party-status ON parties (party_status);
DROP INDEX ix_party_guild ON parties (guild_id);
DROP INDEX ix_party_name ON parties (name);
DROP INDEX ix_party_rank ON parties (party_rank);
DROP INDEX ix_party_size ON parties (maximum_party_size);
DROP TYPE party_status_enum;
DROP TABLE parties;
