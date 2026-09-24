ALTER TABLE survey_responses
    ADD COLUMN tracking_tool         text,
    ADD COLUMN tracking_tool_other   text,
    ADD COLUMN tracking_satisfaction text,
    ADD COLUMN tracking_app_feedback text;
