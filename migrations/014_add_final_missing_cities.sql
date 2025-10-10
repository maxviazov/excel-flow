-- 014_add_final_missing_cities.sql
-- Adding final missing cities found in Excel but not in database

BEGIN;

INSERT OR IGNORE INTO city_codes (city_heb, city_code) VALUES
('נתניה', 'F1012'),
('ירושלים', 'H527'),
('ראשון לציון', 'J1207'),
('חיפה', 'D428'),
('באר שבע', 'N126'),
('פתח תקווה', 'I1137'),
('חולון', 'i400'),
('בני ברק', 'I198');

COMMIT;