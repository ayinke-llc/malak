CREATE TYPE update_status AS ENUM ('draft', 'sent');

CREATE TABLE updates(
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  workspace_id uuid NOT NULL REFERENCES workspaces(id),
  status update_status NOT NULL DEFAULT 'draft',
  reference VARCHAR (220) UNIQUE NOT NULL,
  created_by uuid NOT NULL REFERENCES users(id),
  sent_by uuid NULL REFERENCES users(id),
  content TEXT NOT NULL DEFAULT '',
  metadata jsonb NOT NULL DEFAULT '{}'::jsonb,

  sent_at    TIMESTAMP WITH TIME ZONE,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE
);

ALTER TABLE updates 
  ADD CONSTRAINT updates_reference_check_key 
  CHECK (reference ~ 'update_[a-zA-Z0-9._]+');

CREATE TABLE update_links(
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  reference VARCHAR (220) UNIQUE NOT NULL,
  update_id uuid NOT NULL REFERENCES updates(id),

  expires_at TIMESTAMP WITH TIME ZONE,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE
);

ALTER TABLE update_links 
  ADD CONSTRAINT update_links_reference_check_key 
  CHECK (reference ~ 'link_[a-zA-Z0-9._]+');


CREATE TABLE update_recipients(
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  reference VARCHAR (220) UNIQUE NOT NULL,
  update_id uuid NOT NULL REFERENCES updates(id),
  email VARCHAR (200) UNIQUE NOT NULL,

  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE
);

ALTER TABLE update_recipients 
  ADD CONSTRAINT update_recipients_reference_check_key 
  CHECK (reference ~ 'recipient_[a-zA-Z0-9._]+');


CREATE TYPE update_schedule_status AS ENUM ('scheduled', 'cancelled', 'sent', 'failed');

CREATE TABLE update_schedules(
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  reference VARCHAR (220) UNIQUE NOT NULL,
  update_id uuid NOT NULL REFERENCES updates(id),
  scheduled_by uuid NULL REFERENCES users(id),
  status update_schedule_status NOT NULL DEFAULT 'scheduled',

  send_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE
);

ALTER TABLE update_schedules
  ADD CONSTRAINT update_schedules_reference_check_key 
  CHECK (reference ~ 'schedule_[a-zA-Z0-9._]+');
