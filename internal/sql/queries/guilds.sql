-- name: InsertGuild :one
INSERT INTO guilds (
	name,
	password
) values($1, $2) 
ON CONFLICT DO UPDATE
SET password = excluded.password
RETURNING *;

-- name: GetGuildByID :one
SELECT
id,
current_rank,
treasury
recruitment_slots,
current_rank
FROM guilds
WHERE id =$1;

-- name: GetGuildByName :one
SELECT 
id,
current_rank,
treasury
recruitment_slots,
current_rank,
password
FROM guilds
WHERE name = $1;

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

-- name: UpdateTreasury :exec
UPDATE guilds
SET treasury = treasury + $2
WHERE id = $1
RETURNING *;

-- name: SetGuildRank :exec
UPDATE guilds
SET current_rank = $2
WHERE id = $1
RETURNING *;

-- name: SetRecruitmentSlots :exec
UPDATE guilds
SET recruitment_slots = $2
WHERE id = $1
RETURNING *;
