-- positive_reason/neutral_reason were single-value (text) because those branch
-- questions were single-select; they're now multi-select on the frontend, so
-- these become arrays, matching how negative_reasons (already multi-select)
-- is stored.
ALTER TABLE survey_responses
    DROP COLUMN positive_reason,
    DROP COLUMN neutral_reason,
    ADD COLUMN positive_reasons jsonb,
    ADD COLUMN neutral_reasons jsonb,
    ADD COLUMN negative_reason_other text;
