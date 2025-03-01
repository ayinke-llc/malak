CREATE TABLE system_templates (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  reference VARCHAR (220) UNIQUE NOT NULL,
  content jsonb NOT NULL DEFAULT '{}'::jsonb,
  title VARCHAR(210) UNIQUE NOT NULL,
  description TEXT NOT NULL,
  number_of_uses SMALLINT NOT NULL DEFAULT 0,

  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE
);

ALTER TABLE system_templates ADD CONSTRAINT system_templates_reference_check_key CHECK (reference ~ 'system_templates_[a-zA-Z0-9._]+');
