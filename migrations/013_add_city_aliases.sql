-- 013_add_city_aliases.sql
-- Adding common city name aliases

BEGIN;

-- Добавляем алиасы для сокращенных названий городов
INSERT OR IGNORE INTO city_aliases (alias_heb, target_heb) VALUES ('מעלות', 'מעלות-תרשיחא');
INSERT OR IGNORE INTO city_aliases (alias_heb, target_heb) VALUES ('ת"א', 'תל אביב-יפו');
INSERT OR IGNORE INTO city_aliases (alias_heb, target_heb) VALUES ('תל-אביב', 'תל אביב-יפו');
INSERT OR IGNORE INTO city_aliases (alias_heb, target_heb) VALUES ('ירושלים', 'ירושלים');
INSERT OR IGNORE INTO city_aliases (alias_heb, target_heb) VALUES ('חיפה', 'חיפה');

COMMIT;