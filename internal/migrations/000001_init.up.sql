CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name            text NOT NULL,
    email           text UNIQUE,
    phone           text UNIQUE,
    auth_method     text NOT NULL DEFAULT 'email',
    last_login_at   timestamptz,
    login_count     integer NOT NULL DEFAULT 0,
    created_at      timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT users_email_or_phone CHECK (email IS NOT NULL OR phone IS NOT NULL)
);

CREATE TABLE verticals (
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    key             text NOT NULL UNIQUE,
    label           text NOT NULL,
    icon            text,
    input_schema    jsonb NOT NULL,
    created_at      timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE challenges (
    id                  uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    title               text NOT NULL,
    vertical_id         uuid NOT NULL REFERENCES verticals(id),
    influencer_handle   text,
    member_count        integer NOT NULL DEFAULT 0,
    difficulty_stat     text,
    is_template         boolean NOT NULL DEFAULT true,
    created_at          timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE user_challenges (
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         uuid NOT NULL REFERENCES users(id),
    challenge_id    uuid NOT NULL REFERENCES challenges(id),
    custom_goal     jsonb,
    status          text NOT NULL DEFAULT 'active',
    joined_at       timestamptz NOT NULL DEFAULT now(),
    UNIQUE (user_id, challenge_id)
);

CREATE TABLE checkins (
    id                      uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_challenge_id       uuid NOT NULL REFERENCES user_challenges(id),
    date                    date NOT NULL DEFAULT CURRENT_DATE,
    verification_status     text NOT NULL DEFAULT 'auto_approved',
    metric_data             jsonb NOT NULL,
    created_at              timestamptz NOT NULL DEFAULT now(),
    UNIQUE (user_challenge_id, date)
);

CREATE TABLE streaks (
    user_challenge_id   uuid PRIMARY KEY REFERENCES user_challenges(id),
    current_streak      integer NOT NULL DEFAULT 0,
    longest_streak      integer NOT NULL DEFAULT 0,
    last_checkin_date   date,
    updated_at          timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE wallet_transactions (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     uuid NOT NULL REFERENCES users(id),
    delta       integer NOT NULL,
    reason      text NOT NULL,
    created_at  timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX wallet_transactions_user_id_idx ON wallet_transactions(user_id);

CREATE TABLE redemption_items (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    type        text NOT NULL,
    title       text NOT NULL,
    coin_cost   integer NOT NULL,
    metadata    jsonb NOT NULL DEFAULT '{}'::jsonb,
    active      boolean NOT NULL DEFAULT true,
    created_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE redemptions (
    id                      uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id                 uuid NOT NULL REFERENCES users(id),
    redemption_item_id      uuid NOT NULL REFERENCES redemption_items(id),
    code_or_slot            text,
    created_at              timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE follows (
    follower_id     uuid NOT NULL REFERENCES users(id),
    followee_id     uuid NOT NULL REFERENCES users(id),
    created_at      timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (follower_id, followee_id)
);
