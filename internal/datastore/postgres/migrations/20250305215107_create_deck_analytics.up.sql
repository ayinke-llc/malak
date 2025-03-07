CREATE TABLE deck_daily_engagements (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  deck_id uuid NOT NULL REFERENCES decks(id),
  workspace_id uuid NOT NULL REFERENCES workspaces(id),
  reference VARCHAR (220) UNIQUE NOT NULL,
  engagement_count INT NOT NULL DEFAULT 0,
  engagement_date DATE NOT NULL,

  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,

  UNIQUE(deck_id,workspace_id,engagement_date)
);

CREATE INDEX idx_deck_daily_stats_deck_date ON deck_daily_engagements(deck_id, engagement_date);

ALTER TABLE deck_daily_engagements ADD CONSTRAINT reference_check_key 
  CHECK (reference ~ 'deck_daily_engagement_[a-zA-Z0-9._]+');

CREATE TABLE deck_analytics (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  deck_id uuid NOT NULL REFERENCES decks(id),
  reference VARCHAR (220) UNIQUE NOT NULL,
  total_views INT NOT NULL DEFAULT 0,
  unique_viewers INT NOT NULL DEFAULT 0,
  total_downloads INT NOT NULL DEFAULT 0,

  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE
);

ALTER TABLE deck_analytics ADD CONSTRAINT reference_check_key 
  CHECK (reference ~ 'deck_analytic_[a-zA-Z0-9._]+');

CREATE TABLE deck_viewer_sessions (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  deck_id uuid NOT NULL REFERENCES decks(id),
  reference VARCHAR (220) UNIQUE NOT NULL,
  user_id uuid REFERENCES users(id), -- NULL for anonymous users
  session_id VARCHAR(100) UNIQUE, -- For anonymous users
  device_info VARCHAR(255),
  browser VARCHAR(100),
  os VARCHAR(100),
  ip_address VARCHAR(45),
  country VARCHAR(100),
  city VARCHAR(100),
  time_spent_seconds INT NOT NULL DEFAULT 0,
  viewed_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE
);

ALTER TABLE deck_viewer_sessions ADD CONSTRAINT reference_check_key 
  CHECK (reference ~ 'deck_viewer_session_[a-zA-Z0-9._]+');

CREATE TABLE deck_geographic_stats (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  reference VARCHAR (220) UNIQUE NOT NULL,
  deck_id uuid NOT NULL REFERENCES decks(id),
  country VARCHAR(100) NOT NULL,
  view_count INT NOT NULL DEFAULT 0,
  stat_date DATE NOT NULL,

  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,

  UNIQUE(deck_id, country, stat_date)
);

ALTER TABLE deck_geographic_stats ADD CONSTRAINT reference_check_key 
  CHECK (reference ~ 'deck_geographic_stat_[a-zA-Z0-9._]+');

CREATE INDEX idx_viewer_sessions_deck_date ON deck_viewer_sessions(deck_id, viewed_at);
CREATE INDEX idx_viewer_sessions_user ON deck_viewer_sessions(user_id);
CREATE INDEX idx_geographic_stats_deck_date ON deck_geographic_stats(deck_id, stat_date);
