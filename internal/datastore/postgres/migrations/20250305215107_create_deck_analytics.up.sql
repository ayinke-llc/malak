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
