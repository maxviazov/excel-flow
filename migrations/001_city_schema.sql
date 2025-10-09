-- 001_city_schema.sql
-- Schema for official MOH city codes and alias normalization
-- SQLite 3.38+ recommended

PRAGMA foreign_keys = ON;
PRAGMA journal_mode = WAL;

BEGIN;

-- Источник и контроль версий справочника (какой файл загрузили, когда, хэш)
CREATE TABLE IF NOT EXISTS city_meta (
  id              INTEGER PRIMARY KEY,
  source_name     TEXT NOT NULL,              -- например: FCS-1760033233044.xlsx
  source_sheet    TEXT NOT NULL,              -- Sheet1 (или фактическое имя)
  source_hash     TEXT NOT NULL,              -- SHA256 содержимого при загрузке
  loaded_at       TEXT NOT NULL DEFAULT (datetime('now','localtime'))
);

-- Канонический справочник МОЗ (имена ивритом + официальный код)
CREATE TABLE IF NOT EXISTS city_codes (
  id          INTEGER PRIMARY KEY,
  city_heb    TEXT NOT NULL,                  -- שם יישוב (каноническое имя)
  city_code   TEXT NOT NULL,                  -- קוד יישוב (например: J112)
  active      INTEGER NOT NULL DEFAULT 1,     -- флаг активна/устарела
  updated_at  TEXT NOT NULL DEFAULT (datetime('now','localtime')),
  UNIQUE (city_heb),
  UNIQUE (city_code)
);

-- Алиасы «грязных» названий → каноническое имя
CREATE TABLE IF NOT EXISTS city_aliases (
  alias_heb   TEXT PRIMARY KEY,               -- например: 'ת"א', 'תל-אביב'
  target_heb  TEXT NOT NULL,                  -- должно существовать в city_codes.city_heb
  updated_at  TEXT NOT NULL DEFAULT (datetime('now','localtime')),
  FOREIGN KEY (target_heb) REFERENCES city_codes(city_heb)
    ON UPDATE CASCADE
    ON DELETE RESTRICT
);

-- Ускоряющие индексы
CREATE INDEX IF NOT EXISTS idx_city_codes_city_heb   ON city_codes (city_heb);
CREATE INDEX IF NOT EXISTS idx_city_codes_city_code  ON city_codes (city_code);
CREATE INDEX IF NOT EXISTS idx_city_aliases_target   ON city_aliases (target_heb);

-- Представление «всё про город»: прямое имя или алиас → код
CREATE VIEW IF NOT EXISTS v_city_lookup AS
SELECT c.city_heb      AS key_heb,
       c.city_heb      AS canon_heb,
       c.city_code     AS city_code,
       0               AS is_alias
FROM city_codes c
UNION ALL
SELECT a.alias_heb     AS key_heb,
       a.target_heb    AS canon_heb,
       c.city_code     AS city_code,
       1               AS is_alias
FROM city_aliases a
JOIN city_codes  c ON c.city_heb = a.target_heb;

COMMIT;
