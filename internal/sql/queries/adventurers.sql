-- name: GetRecruitableAdventurers :many
SELECT 
id,
joined_at,
current_rank,
current_activity,
name,
role 
FROM adventurers
WHERE guild_id IS NULL
ORDER BY
	CASE WHEN sqlc.arg(sort_by)::text = 'joined_at' THEN joined_at END ASC,
	CASE WHEN sqlc.arg(sort_by)::text = 'name' THEN name END ASC,
	CASE WHEN sqlc.arg(sort_by)::text = 'role' THEN role END ASC
	CASE WHEN sqlc.arg(sort_by)::text = 'status' THEN role END ASC;


