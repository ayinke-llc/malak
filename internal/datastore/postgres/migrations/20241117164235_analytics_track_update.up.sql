CREATE TABLE update_stats(
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  reference VARCHAR (220) UNIQUE NOT NULL,

  update_id uuid NOT NULL REFERENCES updates(id),
  total_opens INT NOT NULL DEFAULT 0,
  total_reactions INT NOT NULL DEFAULT 0,
  total_clicks INT NOT NULL DEFAULT 0,
  total_sent INT NOT NULL DEFAULT 0,
  unique_opens INT NOT NULL DEFAULT 0,

  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE
);

ALTER TABLE update_stats
  ADD CONSTRAINT update_stats_reference_check_key 
  CHECK (reference ~ 'update_stat_[a-zA-Z0-9._-]+');

CREATE TABLE update_recipient_stats(
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  reference VARCHAR (220) UNIQUE NOT NULL,
  recipient_id uuid NOT NULL REFERENCES update_recipients(id),

  has_reaction BOOLEAN NOT NULL DEFAULT FALSE,
  is_delivered BOOLEAN NOT NULL DEFAULT FALSE,
  is_bounced BOOLEAN NOT NULL DEFAULT FALSE,
  last_opened_at TIMESTAMP WITH TIME ZONE,

  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE
);

ALTER TABLE update_recipient_stats 
  ADD CONSTRAINT update_recipient_stats_check_key 
  CHECK (reference ~ 'recipient_stat_[a-zA-Z0-9._-]+');

