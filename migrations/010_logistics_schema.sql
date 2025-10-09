-- 010_logistics_schema.sql
-- Schema for logistics data (drivers, vehicles, routes)

PRAGMA foreign_keys = ON;
PRAGMA journal_mode = WAL;

BEGIN;

-- Metadata table
CREATE TABLE IF NOT EXISTS logistics_meta (
  id              INTEGER PRIMARY KEY,
  source_name     TEXT NOT NULL,
  loaded_at       TEXT NOT NULL DEFAULT (datetime('now','localtime'))
);

-- Drivers table
CREATE TABLE IF NOT EXISTS drivers (
  id              INTEGER PRIMARY KEY,
  driver_id       TEXT NOT NULL UNIQUE,
  driver_name     TEXT NOT NULL,
  license_number  TEXT,
  phone           TEXT,
  active          INTEGER NOT NULL DEFAULT 1,
  created_at      TEXT NOT NULL DEFAULT (datetime('now','localtime'))
);

-- Vehicles table
CREATE TABLE IF NOT EXISTS vehicles (
  id              INTEGER PRIMARY KEY,
  vehicle_id      TEXT NOT NULL UNIQUE,
  license_plate   TEXT NOT NULL,
  vehicle_type    TEXT,
  capacity_kg     INTEGER,
  active          INTEGER NOT NULL DEFAULT 1,
  created_at      TEXT NOT NULL DEFAULT (datetime('now','localtime'))
);

-- Routes table
CREATE TABLE IF NOT EXISTS routes (
  id              INTEGER PRIMARY KEY,
  route_id        TEXT NOT NULL UNIQUE,
  route_name      TEXT NOT NULL,
  start_location  TEXT,
  end_location    TEXT,
  active          INTEGER NOT NULL DEFAULT 1,
  created_at      TEXT NOT NULL DEFAULT (datetime('now','localtime'))
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_drivers_driver_id ON drivers (driver_id);
CREATE INDEX IF NOT EXISTS idx_vehicles_vehicle_id ON vehicles (vehicle_id);
CREATE INDEX IF NOT EXISTS idx_routes_route_id ON routes (route_id);

COMMIT;