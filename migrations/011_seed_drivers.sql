-- 011_seed_drivers.sql
-- Seed data for logistics

BEGIN;

-- Insert metadata
INSERT INTO logistics_meta (source_name) VALUES ('initial_seed');

-- Insert sample drivers
INSERT OR IGNORE INTO drivers (driver_id, driver_name, license_number, phone) VALUES
('DRV001', 'יוסי כהן', '12345678', '050-1234567'),
('DRV002', 'משה לוי', '23456789', '052-2345678'),
('DRV003', 'דוד אברהם', '34567890', '054-3456789'),
('DRV004', 'אבי יצחק', '45678901', '053-4567890'),
('DRV005', 'רון יעקב', '56789012', '050-5678901');

-- Insert sample vehicles
INSERT OR IGNORE INTO vehicles (vehicle_id, license_plate, vehicle_type, capacity_kg) VALUES
('VEH001', '123-45-678', 'משאית', 5000),
('VEH002', '234-56-789', 'משאית', 7500),
('VEH003', '345-67-890', 'ואן', 2000),
('VEH004', '456-78-901', 'משאית', 10000),
('VEH005', '567-89-012', 'ואן', 1500);

-- Insert sample routes
INSERT OR IGNORE INTO routes (route_id, route_name, start_location, end_location) VALUES
('RT001', 'תל אביב - ירושלים', 'תל אביב', 'ירושלים'),
('RT002', 'חיפה - באר שבע', 'חיפה', 'באר שבע'),
('RT003', 'נתניה - אשדוד', 'נתניה', 'אשדוד'),
('RT004', 'פתח תקווה - אילת', 'פתח תקווה', 'אילת'),
('RT005', 'רמת גן - חולון', 'רמת גן', 'חולון');

COMMIT;