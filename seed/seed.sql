-- Seed data for local/dev testing. Idempotent-ish via ON CONFLICT DO NOTHING on natural keys.

INSERT INTO verticals (key, label, icon, input_schema) VALUES
('reading', 'Reading', '📖', '[{"key":"pages","type":"number","label":"Pages read","unit":"pages"}]'),
('fitness', 'Fitness / Gym', '🏋️', '[{"key":"workout_name","type":"text","label":"Workout"},{"key":"exercises","type":"exercise_list","label":"Exercises"}]'),
('weight', 'Weight', '⚖️', '[{"key":"weight_kg","type":"number","label":"Weight","unit":"kg"}]'),
('diet', 'Diet / Food', '🥗', '[{"key":"meal","type":"select","label":"Meal","options":["breakfast","lunch","dinner","snack"]},{"key":"photo","type":"photo","label":"Photo"}]')
ON CONFLICT (key) DO NOTHING;

-- disqualify_after_missed_days: only set on challenges whose own framing implies
-- real stakes (Reading, Fitness). Left NULL for Weight/Diet — auto-eliminating
-- someone over a body-related measure runs against the product's anti-guilt
-- stance, so those stay non-punitive until a real Freeze Days mechanic exists.
INSERT INTO challenges (title, vertical_id, influencer_handle, difficulty_stat, is_template, disqualify_after_missed_days)
SELECT '75-Day Reading Streak', id, '@bookinfluencer', 'only 8% finish', true, 3 FROM verticals WHERE key = 'reading'
UNION ALL
SELECT '30-Day Gym Streak', id, '@fitwithrahul', 'top 3% get free coaching', true, 2 FROM verticals WHERE key = 'fitness'
UNION ALL
SELECT '12-Week Weight Challenge', id, '@fitwithrahul', 'avg loss 4.2kg', true, NULL FROM verticals WHERE key = 'weight'
UNION ALL
SELECT '21-Day Clean Eating', id, '@nutribyneha', 'new this week', true, NULL FROM verticals WHERE key = 'diet';

INSERT INTO redemption_items (type, title, coin_cost, metadata, active) VALUES
('voucher', '20% off — Decathlon', 300, '{"brand":"Decathlon"}', true),
('voucher', 'Free sample — MuscleBlaze', 150, '{"brand":"MuscleBlaze"}', true),
('voucher', '₹200 off — Cult.fit gear', 400, '{"brand":"Cult.fit"}', true),
('voucher', '₹150 off — Amazon', 350, '{"brand":"Amazon"}', true),
('consultation', 'Riya Sharma — Certified Dietician', 1200, '{"duration_min":30,"rating":4.9}', true),
('consultation', 'Karan Mehta — Strength Coach', 1000, '{"duration_min":30,"rating":4.8}', true);
