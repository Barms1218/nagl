-- name: GetRecruitableAdventurers :many
SELECT 
id,
name,
role,
description,
current_rank
FROM adventurers
WHERE guild_id IS NULL
ORDER BY
	CASE WHEN sqlc.arg(sort_by)::text = 'joined_at' THEN joined_at END ASC,
	CASE WHEN sqlc.arg(sort_by)::text = 'name' THEN name END ASC,
	CASE WHEN sqlc.arg(sort_by)::text = 'role' THEN role END ASC,
	CASE WHEN sqlc.arg(sort_by)::text = 'status' THEN role END ASC;


-- name: GetAdventurersByGuild :many
SELECT
id,
name,
current_rank,
role,
description,
current_activity
FROM adventurers
WHERE guild_id = $1;

-- name: GetAdventurerDetails :one
SELECT
id,
name,
current_rank,
role,
current_activity,
recruitment_cost
FROM adventurers
WHERE id = $1;

-- name: GetAdventurers :many
SELECT 
id,
joined_at,
current_rank,
current_activity,
description,
name,
role
FROM adventurers
WHERE guild_id = $1
ORDER BY
	CASE WHEN sqlc.arg(sort_by)::text = 'joined_at' THEN joined_at END ASC,
	CASE WHEN sqlc.arg(sort_by)::text = 'name' THEN name END ASC,
	CASE WHEN sqlc.arg(sort_by)::text = 'role' THEN role END ASC,
	CASE WHEN sqlc.arg(sort_by)::text = 'activity' THEN current_activity END ASC;

-- name: GetAdventurersWithStatus :many
SELECT 
id,
joined_at,
current_rank,
current_activity,
description,
name,
role
FROM adventurers 
WHERE guild_id = $1 AND current_activity = $2
ORDER BY
	CASE WHEN sqlc.arg(sort_by)::text = 'joined_at' THEN joined_at END ASC,
	CASE WHEN sqlc.arg(sort_by)::text = 'name' THEN name END ASC,
	CASE WHEN sqlc.arg(sort_by)::text = 'role' THEN role END ASC;

-- name: GetAdventurersWithRole :many
SELECT 
id,
joined_at,
current_rank,
current_activity,
name,
description,
role 
FROM adventurers
WHERE guild_id = $1 AND role = $2
ORDER BY
	CASE WHEN sqlc.arg(sort_by)::text = 'joined_at' THEN joined_at END ASC,
	CASE WHEN sqlc.arg(sort_by)::text = 'name' THEN name END ASC,
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

-- name: GetAdventurerUpkeepCost :exec
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
SET guild_id = $1
WHERE id = $2;

-- name: SetAdventurerActivity :exec
UPDATE adventurers
SET current_activity = $1
WHERE id = $2;

-- name: SetAdventurerRank :exec
UPDATE adventurers
SET current_rank = $1
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
