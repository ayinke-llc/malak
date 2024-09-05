ALTER TABLE roles DROP COLUMN role;
CREATE TYPE user_role AS ENUM ('admin', 'member', 'billing', 'investor', 'guest');
ALTER TABLE roles ADD COLUMN role user_role NOT NULL DEFAULT 'admin';
