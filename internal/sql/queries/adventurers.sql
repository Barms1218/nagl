-- name: ListRecruitableAdventurers :many
SELECT 
id,
name,
role,
current_rank
FROM adventurers
WHERE guild_id IS NULL
AND (sqlc.narg('name')::text IS NULL OR name = sqlc.narg('name'))
AND (sqlc.narg('role')::role_enum IS NULL OR role = sql.narg('role'))
AND (sqlc.narg('max_rank')::int IS NULL OR current_rank <= sqlc.narg('max_rank'))
AND (sqlc.narg('min_rank')::int IS NULL OR current_rank >= sqlc.narg('min_rank'))
ORDER BY
	CASE WHEN sqlc.arg(sort_by)::text = 'joined_at' THEN joined_at END ASC,
	CASE WHEN sqlc.arg(sort_by)::text = 'name' THEN name END ASC,
	CASE WHEN sqlc.arg(sort_by)::text = 'role' THEN role END ASC,
	CASE WHEN sqlc.arg(sort_by)::text = 'status' THEN role END ASC;


-- name: ListGuildMembers :many
SELECT
id,
name,
current_rank,
role,
current_activity
FROM adventurers
WHERE guild_id = $1
AND (sqlc.narg('name')::text IS NULL OR name = sqlc.narg('name'))
AND (sqlc.narg('role')::role_enum IS NULL OR role = sql.narg('role'))
AND (sqlc.narg('max_rank')::int IS NULL OR current_rank <= sqlc.narg('max_rank'))
AND (sqlc.narg('min_rank')::int IS NULL OR current_rank >= sqlc.narg('min_rank'))
AND (sqlc.narg('current_activity')::activity_enum IS NULL OR current_activity = sqlc.narg('current_activity'))
ORDER BY
	CASE WHEN sqlc.arg(sort_by)::text = 'joined_at' THEN joined_at END ASC,
	CASE WHEN sqlc.arg(sort_by)::text = 'name' THEN name END ASC,
	CASE WHEN sqlc.arg(sort_by)::text = 'role' THEN role END ASC,
	CASE WHEN sqlc.arg(sort_by)::text = 'status' THEN role END ASC;

-- name: GetAdventurerDetails :one
SELECT
id,
party_id,
name,
current_rank,
role,
current_activity,
upkeep_cost
FROM adventurers
WHERE id = $1;

-- name: GetGuildMembers :many
SELECT 
id,
current_rank,
current_activity,
name,
role
FROM adventurers
WHERE guild_id = $1
ORDER BY
	CASE WHEN sqlc.arg(sort_by)::text = 'joined_at' THEN joined_at END ASC,
	CASE WHEN sqlc.arg(sort_by)::text = 'name' THEN name END ASC,
	CASE WHEN sqlc.arg(sort_by)::text = 'role' THEN role END ASC,
	CASE WHEN sqlc.arg(sort_by)::text = 'activity' THEN current_activity END ASC;


-- name: GetAdventurerActivities :many
SELECT 
ah.id,
ah.occurred_at,
ah.activity,
a.name
FROM adventurer_history ah
JOIN adventurers a ON ah.adventurer_id = a.id
WHERE a.id = $1
ORDER BY ah.occurred_at DESC;

-- name: GetAdventurerContractHistory :many
SELECT
a.id AS adventurer_id,
a.name,
c.id AS contract_id,
c.title,
c.contract_status
FROM adventurer_contract_history ach
JOIN adventurers a ON ach.adventurer_id = a.id
JOIN contracts c ON ach.contract_id = c.id
WHERE a.id = $1
ORDER BY ach.occurred_at DESC;

-- name: GetAdventurerUpkeepCost :one
SELECT
upkeep_cost
FROM adventurers
WHERE id = $1;

-- name: UpsertAdventurer :one
Insert INTO adventurers(name, current_rank, role, description)
VALUES($1, $2, $3, $4)
RETURNING *;

-- name: SetAdventurerHired :exec
UPDATE adventurers
SET guild_id = $1, updated_at = CURRENT_TIMESTAMP
WHERE id = $2;

-- name: SetAdventurerActivity :exec
UPDATE adventurers
SET current_activity = $1, updated_at = CURRENT_TIMESTAMP
WHERE id = $2;

-- name: SetAdventurerRank :exec
UPDATE adventurers
SET current_rank = $1, updated_at = CURRENT_TIMESTAMP
WHERE id = $2;

-- name: InsertAdventurerHistory :exec
INSERT INTO adventurer_history(
adventurer_id,
activity
)VALUES($1, $2);

-- name: InsertAdventurerContractHistory :exec
INSERT INTO adventurer_contract_history(
	adventurer_id,
	contract_id,
	status
) VALUES($1, $2, $3);
