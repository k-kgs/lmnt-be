-- Participation lifecycle: leave/disqualify tracking + configurable
-- auto-disqualification threshold per challenge. member_count moves from a
-- static stored column to a live-computed value (see queries/config.sql),
-- so it's dropped here rather than left around as a misleading duplicate.

ALTER TABLE user_challenges
    ADD COLUMN left_at         timestamptz,
    ADD COLUMN disqualified_at timestamptz;

ALTER TABLE challenges
    ADD COLUMN disqualify_after_missed_days integer;

ALTER TABLE challenges
    DROP COLUMN member_count;
