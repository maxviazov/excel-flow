-- 017_add_spelling_aliases.sql
-- Adding aliases for different spelling variants

BEGIN;

-- Алиасы для разных вариантов написания
INSERT OR IGNORE INTO city_aliases (alias_heb, target_heb) VALUES 
('קירית ביאליק', 'קריית ביאליק'),
('קרית ביאליק', 'קריית ביאליק'),
('נהריה', 'נהרייה');

COMMIT;