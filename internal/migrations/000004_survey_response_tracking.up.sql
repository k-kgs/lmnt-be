ALTER TABLE survey_responses
    ADD COLUMN client_id   text,
    ADD COLUMN completed   boolean NOT NULL DEFAULT false,
    ADD COLUMN last_screen text,
    ADD COLUMN updated_at  timestamptz NOT NULL DEFAULT now();

-- No rows exist yet in practice (000003 was never deployed before this migration
-- was written), but backfill defensively so the NOT NULL + UNIQUE below are safe
-- to add regardless.
UPDATE survey_responses SET client_id = id::text WHERE client_id IS NULL;

ALTER TABLE survey_responses ALTER COLUMN client_id SET NOT NULL;
ALTER TABLE survey_responses ADD CONSTRAINT survey_responses_client_id_key UNIQUE (client_id);
