CREATE TABLE contacts (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  email VARCHAR (220) NOT NULL,
  workspace_id uuid NOT NULL REFERENCES workspaces(id),
  reference VARCHAR (220) UNIQUE NOT NULL,
  first_name VARCHAR (200) NOT NULL DEFAULT '',
  last_name VARCHAR (200) NOT NULL DEFAULT '',
  company VARCHAR (220) NOT NULL DEFAULT '' ,
  city VARCHAR (220) NOT NULL DEFAULT '' ,
  phone VARCHAR (50) NOT NULL DEFAULT '',
  notes TEXT NOT NULL DEFAULT '',
  owner_id uuid NOT NULL REFERENCES users(id),
  created_by uuid NOT NULL REFERENCES users(id),

  metadata jsonb NOT NULL DEFAULT '{}'::jsonb,

  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE
);

ALTER TABLE contacts ADD CONSTRAINT contacts_reference_check_key CHECK (reference ~ 'contact_[a-zA-Z0-9._]+');
