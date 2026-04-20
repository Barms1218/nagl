-- +goose Up
CREATE TABLE guilds (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name TEXT UNIQUE NOT NULL,
	password TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	recruitment_slots INTEGER NOT NULL DEFAULT 3,
	treasury INTEGER NOT NULL DEFAULT 100,
	current_rank INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX ix_guild_name ON guilds (name);
CREATE INDEX ix_guild_rank ON guilds (current_rank);

-- +goose Down
DROP INDEX ix_guild_name ON guilds (name);
DROP INDEX ix_guild_rank ON guilds (current_rank);
DROP TABLE guilds;
