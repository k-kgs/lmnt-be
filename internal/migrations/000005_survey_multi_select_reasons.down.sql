ALTER TABLE survey_responses
    DROP COLUMN positive_reasons,
    DROP COLUMN neutral_reasons,
    DROP COLUMN negative_reason_other,
    ADD COLUMN positive_reason text,
    ADD COLUMN neutral_reason text;
