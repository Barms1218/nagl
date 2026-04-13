-- +goose Up
CREATE TYPE contract_status_enum AS ENUM(
	'available',
	'in-progress',
	'complete', 
	'failed'
);

CREATE TABLE contracts (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	guild_id UUID NOT NULL REFERENCES guilds(id) ON DELETE CASCADE,
	title TEXT,
	difficulty INTEGER NOT NULL DEFAULT 1 ,
	minimum_party_size INTEGER NOT NULL DEFAULT 1,
	CHECK (difficulty >= 1 AND difficulty <= 5),
	contract_status contract_status_enum NOT NULL DEFAULT 'available',
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TYPE contract_status;
DROP TABLE quests;
