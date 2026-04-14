-- name: InsertGuild :one
INSERT INTO guilds (
	name 
) values($1) 
ON CONFLICT DO UPDATE
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
SELECT treasury
FROM guilds
WHERE id = $1;

-- name: GetRecruitmentSlots :one
SELECT recruitment_slots
FROM guilds
WHERE id = $1;

-- name: GetCurrentRank :one
SELECT current_rank
FROM guilds
WHERE id = $1;

-- name: UpdateTreasury :one
UPDATE guilds
SET treasury = treasury + $2
WHERE id = $1
RETURNING *;

-- name: SetGuildRank :one
UPDATE guilds
SET current_rank = $2
WHERE id = $1
RETURNING *;

-- name: SetRecruitmentSlots :one
UPDATE guilds
SET recruitment_slots = $2
WHERE id = $1
RETURNING *;
