-- 016_allow_multiple_codes_per_city.sql
-- Allowing multiple codes per city by removing UNIQUE constraint on city_heb

BEGIN;

-- Создаем новую таблицу без ограничения уникальности на city_heb
CREATE TABLE city_codes_new (
  id          INTEGER PRIMARY KEY,
  city_heb    TEXT NOT NULL,
  city_code   TEXT NOT NULL,
  active      INTEGER NOT NULL DEFAULT 1,
  updated_at  TEXT NOT NULL DEFAULT (datetime('now','localtime')),
  UNIQUE (city_code)  -- только код должен быть уникальным
);

-- Копируем данные
INSERT INTO city_codes_new (id, city_heb, city_code, active, updated_at)
SELECT id, city_heb, city_code, active, updated_at FROM city_codes;

-- Удаляем старую таблицу и переименовываем новую
DROP TABLE city_codes;
ALTER TABLE city_codes_new RENAME TO city_codes;

-- Восстанавливаем индексы
CREATE INDEX idx_city_codes_city_heb   ON city_codes (city_heb);
CREATE INDEX idx_city_codes_city_code  ON city_codes (city_code);

-- Теперь добавляем альтернативные коды
INSERT INTO city_codes (city_heb, city_code) VALUES
('נתניה', 'F1012'),
('ירושלים', 'H527'), 
('ראשון לציון', 'J1207'),
('חיפה', 'D428'),
('באר שבע', 'N126'),
('פתח תקווה', 'I1137'),
('חולון', 'i400'),
('בני ברק', 'I198');

COMMIT;