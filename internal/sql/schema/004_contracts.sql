-- +goose Up
CREATE TYPE contract_status AS ENUM('available', 'in-progress', 'complete', 'failed')
CREATE TABLE contracts (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	contract_id UUID NOT NULL REFERENCES contracts(id) ON DELETE CASCADE,
	title TEXT,
	difficulty INTEGER NOT NULL DEFAULT 1 ,
	minimum_party_size INTEGER NOT NULL DEFAULT 1,
	reward NOT NULL DEFAULT (random(25, 75) * difficulty,
	CHECK (difficulty >= 1 AND difficulty <= 6),
	status contract_status NOT NULL DEFAULT 'available',
)

-- +goose Down
DROP TABLE quests;
