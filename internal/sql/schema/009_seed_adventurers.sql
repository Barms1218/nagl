-- +goose Up
INSERT INTO adventurers(
	id,
	guild_id,
	name,
	current_rank,
	current_activity,
	description,
	role,
	upkeep_cost
)VALUES(
gen_random_uuid(),

)
