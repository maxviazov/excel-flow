-- 002_seed_city_codes.sql
-- Seed data for city codes and aliases

BEGIN;

-- Insert meta information
INSERT OR REPLACE INTO city_meta (source_name, source_sheet, source_hash)
VALUES ('seed_data', 'manual', 'initial_seed');

-- Insert canonical city codes
INSERT OR REPLACE INTO city_codes (city_heb, city_code) VALUES
('תל אביב-יפו', 'J112'),
('ירושלים', 'J001'),
('חיפה', 'J201'),
('ראשון לציון', 'J301'),
('פתח תקווה', 'J401'),
('אשדוד', 'J501'),
('נתניה', 'J601'),
('באר שבע', 'J701'),
('בני ברק', 'J801'),
('חולון', 'J901');

-- Insert common aliases
INSERT OR REPLACE INTO city_aliases (alias_heb, target_heb) VALUES
('ת"א', 'תל אביב-יפו'),
('תל-אביב', 'תל אביב-יפו'),
('תל אביב', 'תל אביב-יפו'),
('יפו', 'תל אביב-יפו'),
('ירושלים הבירה', 'ירושלים'),
('חיפה העיר', 'חיפה'),
('ראשל"צ', 'ראשון לציון'),
('פ"ת', 'פתח תקווה'),
('פתח-תקווה', 'פתח תקווה'),
('ב"ש', 'באר שבע'),
('באר-שבע', 'באר שבע');

COMMIT;