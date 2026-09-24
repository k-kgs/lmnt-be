CREATE TABLE survey_responses (
    id                  uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    track               text,
    track_other         text,
    pivot_importance    text,
    pivot_satisfaction  text,
    branch              text,
    positive_reason     text,
    neutral_reason      text,
    negative_reasons    jsonb,
    reward_kano         text,
    monetization        text,
    age                 text,
    gender              text,
    email               text,
    email_choice        text,
    persona_key         text,
    created_at          timestamptz NOT NULL DEFAULT now()
);
