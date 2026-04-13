-- name: InsertGuild: one
INSERT INTO guilds (
	name, 
) values($1) 
ON CONFLICT(name) DO UPDATE
SET name = excluded.name
RETURNING *;

-- name: GetGuildMembers :many
SELECT 
id,
joined_at,
current_rank,
current_activity,
name,
role
FROM adventurers
WHERE guild_id = $1
ORDER BY
	CASE WHEN sqlc.arg(sort_by)::text = 'joined_at' THEN joined_at END ASC,
	CASE WHEN sqlc.arg(sort_by)::text = 'name' THEN name END ASC,
	CASE WHEN sqlc.arg(sort_by)::text = 'role' THEN role END ASC
	CASE WHEN sqlc.arg(sort_by)::text = 'status' THEN role END ASC;

-- name: GetMembersWithStatus: many
SELECT 
id,
joined_at,
current_rank,
current_activity,
name,
role
FROM adventurers 
WHERE guild_id = $1 AND current_activity = $2
ORDER BY
	CASE WHEN sqlc.arg(sort_by)::text = 'joined_at' THEN joined_at END ASC,
	CASE WHEN sqlc.arg(sort_by)::text = 'name' THEN name END ASC
	CASE WHEN sqlc.arg(sort_by)::text = 'role' THEN role END ASC;

-- name: GetMembersWithRole :many
SELECT 
id,
joined_at,
current_rank,
current_activity,
name,
role 
FROM adventurers
WHERE guild_id = $1 AND role = $2
ORDER BY
	CASE WHEN sqlc.arg(sort_by)::text = 'joined_at' THEN joined_at END ASC,
	CASE WHEN sqlc.arg(sort_by)::text = 'name' THEN name END ASC
	CASE WHEN sqlc.arg(sort_by)::text = 'status' THEN role END ASC;


