-- 015_add_alternative_city_codes.sql
-- Adding alternative city codes for cities that have multiple codes in MOH list

BEGIN;

-- Временно отключаем ограничение уникальности для city_heb
-- Добавляем альтернативные коды для городов
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