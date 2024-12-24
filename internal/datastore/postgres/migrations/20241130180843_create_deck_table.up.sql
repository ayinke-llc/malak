CREATE TABLE decks (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  reference VARCHAR (220) UNIQUE NOT NULL,
  workspace_id uuid NOT NULL REFERENCES workspaces(id),
  created_by uuid NOT NULL REFERENCES users(id),
  title VARCHAR (220) NOT NULL,
  short_link VARCHAR (50) NOT NULL,

  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE
);

ALTER TABLE decks ADD CONSTRAINT deck_reference_check_key CHECK (reference ~ 'deck_[a-zA-Z0-9._]+');

CREATE TABLE deck_preferences (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  reference VARCHAR (220) UNIQUE NOT NULL,
  workspace_id uuid NOT NULL REFERENCES workspaces(id),
  deck_id uuid NOT NULL REFERENCES decks(id),
  enable_downloading BOOLEAN NOT NULL DEFAULT false,
  require_email BOOLEAN NOT NULL DEFAULT false,
  password jsonb NOT NULL DEFAULT '{}'::jsonb,
  expires_at TIMESTAMP WITH TIME ZONE,
  created_by uuid NOT NULL REFERENCES users(id),

  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE
);

ALTER TABLE deck_preferences 
  ADD CONSTRAINT deck_preference_reference_check_key 
  CHECK (reference ~ 'deck_preference_[a-zA-Z0-9._]+');
