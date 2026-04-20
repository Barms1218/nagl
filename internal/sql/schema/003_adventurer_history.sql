-- +goose Up
CREATE TABLE adventurer_history (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	adventurer_id UUID NOT NULL REFERENCES adventurers(id) ON DELETE CASCADE,
	party_id UUID REFERENCES parties(id),
	occurred_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	activity activity_enum NOT NULL
);

CREATE INDEX ix_ad_id ON adventurer_history (adventurer_id);
CREATE INDEX ix_party_id ON adventurer_history (party_id);
CREATE INDEX ix_activity ON adventurer_history (activity);

-- +goose Down
DROP INDEX ix_ad_id ON adventurer_history (adventurer_id);
DROP INDEX ix_party_id ON adventurer_history (party_id);
DROP INDEX ix_activity ON adventurer_history (activity);
DROP TABLE adventurer_history;
