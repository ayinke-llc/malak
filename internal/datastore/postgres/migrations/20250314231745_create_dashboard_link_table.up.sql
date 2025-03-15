CREATE TYPE dashboard_link_type AS ENUM('default', 'contact');

CREATE TABLE dashboard_links (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  reference VARCHAR (220) UNIQUE NOT NULL,
  dashboard_id uuid NOT NULL REFERENCES dashboards(id),
  link_type dashboard_link_type NOT NULL,
  token VARCHAR (220) UNIQUE NOT NULL,
  contact_id uuid REFERENCES contacts(id), -- NULL for anonymous users
  expires_at TIMESTAMP WITH TIME ZONE,

  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE
);

ALTER TABLE dashboard_links ADD CONSTRAINT reference_check_key 
  CHECK (reference ~ 'dashboard_link_[a-zA-Z0-9._]+');

-- CREATE TABLE dashboard_link_access_logs (
--   id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
--   reference VARCHAR (220) UNIQUE NOT NULL,
--   dashboard_link_id uuid NOT NULL REFERENCES dashboard_links(id),
--   link_type dashboard_link_type NOT NULL,
--   contact_id uuid REFERENCES contacts(id), -- NULL for anonymous users
--   expires_at TIMESTAMP WITH TIME ZONE,
--
--   created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
--   updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
--   deleted_at TIMESTAMP WITH TIME ZONE
-- );
--
-- ALTER TABLE dashboard_links ADD CONSTRAINT reference_check_key 
--   CHECK (reference ~ 'dashboard_link_[a-zA-Z0-9._]+');
--
