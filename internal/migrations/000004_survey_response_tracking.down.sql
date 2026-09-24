ALTER TABLE survey_responses
    DROP CONSTRAINT survey_responses_client_id_key,
    DROP COLUMN client_id,
    DROP COLUMN completed,
    DROP COLUMN last_screen,
    DROP COLUMN updated_at;
