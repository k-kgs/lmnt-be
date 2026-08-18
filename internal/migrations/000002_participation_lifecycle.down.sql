ALTER TABLE challenges
    ADD COLUMN member_count integer NOT NULL DEFAULT 0;

ALTER TABLE challenges
    DROP COLUMN disqualify_after_missed_days;

ALTER TABLE user_challenges
    DROP COLUMN left_at,
    DROP COLUMN disqualified_at;
