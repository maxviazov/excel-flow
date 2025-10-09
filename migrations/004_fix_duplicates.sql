-- 004_fix_duplicates.sql
-- Fix duplicate city codes

BEGIN;

-- Delete all existing data to start fresh
DELETE FROM city_aliases;
DELETE FROM city_codes;
DELETE FROM city_meta;

-- Insert meta information
INSERT INTO city_meta (source_name, source_sheet, source_hash)
VALUES ('FCS-1760033233044.xlsx', 'ALL_FIXED', 'd535fd1052716a37cdde7c0ddb445ce8bc84860145ad24f0355516b76450d760');

-- Insert city codes with ON CONFLICT DO NOTHING to avoid duplicates
INSERT OR IGNORE INTO city_codes (city_heb, city_code) VALUES
('אשדוד', 'J112'),
('חולון', 'I400');

COMMIT;