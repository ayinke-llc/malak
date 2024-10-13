CREATE TYPE update_type_enum AS ENUM ('preview', 'live');

ALTER TABLE update_schedules ADD COLUMN update_type update_type_enum NOT NULL;
