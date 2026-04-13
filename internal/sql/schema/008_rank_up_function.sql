-- +goose Up
CREATE FUNCTION rank_up()
RETURNS TRIGGER AS $$
DECLARE 
	num_completed INTEGER;
	new_rank INTEGER;
	current_treasury INTEGER;
BEGIN
	SELECT COUNT(*) INTO num_completed
	FROM contract_history
	WHERE status = 'complete' AND guild_id = NEW.guild_id;

	SELECT treasury INTO current_treasury
	FROM guilds
	WHERE id = NEW.guild_id;

	CASE -- Get the rank from number of completed quests
		WHEN num_completed >= 40 AND new_treasury >= 5000 THEN new_rank := 5
		WHEN num_completed >= 30 AND new_treasury >= 4000 THEN new_rank := 4
		WHEN num_completed >= 20 AND new_treasury >= 2000 THEN new_rank := 3
		WHEN num_completed >= 10 AND new_treasury >= 1000 THEN new_rank := 2
		ELSE new_rank := 1;
	END CASE;

	UPDATE guilds
	SET current_rank = new_rank
	where id = NEW.guild_id
	AND current_rank < new_rank;

	RETURN NEW; -- Signals to the trigger that the original insert should process
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER check_rank_up
AFTER INSERT ON contract_history
FOR EACH ROW
WHEN (NEW.status = 'complete')
EXECUTE FUNCTION rank_up();

-- +goose Down
DROP FUNCTION rank_up();
DROP TRIGGER check_rank_up;
