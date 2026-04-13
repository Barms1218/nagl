-- +goose Up
CREATE TABLE adventurer_history (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	adventurer_id UUID NOT NULL REFERENCES adventurers(id) ON DELETE CASCADE,
	occurred_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	activity activity_enum NOT NULL
);

-- +goose Down
DROP TABLE adventurer_history;
