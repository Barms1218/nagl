-- +goose Up
CREATE TYPE contract_status_enum AS ENUM(
	'available',
	'in-progress',
	'complete', 
	'failed'
);

CREATE TABLE contracts (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	guild_id UUID  REFERENCES guilds(id) ON DELETE SET NULL,
	title TEXT,
	difficulty INTEGER NOT NULL DEFAULT 1 ,
	rec_party_size INTEGER NOT NULL DEFAULT 1,
	description TEXT,
	CHECK (difficulty >= 1 AND difficulty <= 5),
	contract_status contract_status_enum NOT NULL DEFAULT 'available',
	reward INTEGER NOT NULL DEFAULT 0,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX ix_contract_title ON contracts (title);
CREATE INDEX ix_contract_diff ON contracts (difficulty);
CREATE INDEX ix_contract_size ON contracts (rec_party_size);
CREATE INDEX ix_contract_size ON contracts (contract_status);
CREATE INDEX ix_contract_reward ON contracts (reward);

-- +goose Down
DROP INDEX ix_contract_title ON contracts (title);
DROP INDEX ix_contract_diff ON contracts (difficulty);
DROP INDEX ix_contract_size ON contracts (rec_party_size);
DROP INDEX ix_contract_size ON contracts (contract_status);
DROP INDEX ix_contract_reward ON contracts (reward);
DROP TYPE contract_status_enum;
DROP TABLE contracts;
