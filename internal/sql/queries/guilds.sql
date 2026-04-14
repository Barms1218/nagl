-- name: InsertGuild: one
INSERT INTO guilds (
	name, 
) values($1) 
ON CONFLICT(name) DO UPDATE
SET name = excluded.name
RETURNING *;

-- name: GetGuild :one
SELECT
id,
current_rank,
treasury
FROM guilds
WHERE id =$1;

-- name: GetGuildTreasury :one
SELECT id, treasury
FROM guilds
WHERE id = $1


