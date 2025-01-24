CREATE TYPE contact_share_type AS ENUM ('update', 'deck', 'dashboard');

CREATE TABLE dashboards(
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  reference VARCHAR (220) UNIQUE NOT NULL
);

ALTER TABLE dashboards ADD CONSTRAINT contact_dashboard_reference_check_key CHECK (reference ~ 'dashboard_[a-zA-Z0-9._]+');

CREATE TABLE contact_shares(
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  reference VARCHAR (220) UNIQUE NOT NULL,
  shared_by uuid NOT NULL REFERENCES users(id),
  contact_id uuid NOT NULL REFERENCES contacts(id),
  item_type contact_share_type NOT NULL,
  item_id uuid NOT NULL,

  shared_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE
);

ALTER TABLE contact_shares ADD CONSTRAINT contact_share_reference_check_key CHECK (reference ~ 'contact_share_[a-zA-Z0-9._]+');

CREATE OR REPLACE FUNCTION validate_item_id()
RETURNS TRIGGER AS $$
BEGIN
    IF (NEW.item_type = 'update') THEN
        IF NOT EXISTS (SELECT 1 FROM updates WHERE id = NEW.item_id) THEN
            RAISE EXCEPTION 'Invalid item_id: No matching record in updates for item_type = update';
        END IF;
    ELSIF (NEW.item_type = 'deck') THEN
        IF NOT EXISTS (SELECT 1 FROM decks WHERE id = NEW.item_id) THEN
            RAISE EXCEPTION 'Invalid item_id: No matching record in decks for item_type = deck';
        END IF;
    ELSIF (NEW.item_type = 'dashboard') THEN
        IF NOT EXISTS (SELECT 1 FROM dashboards WHERE id = NEW.item_id) THEN
            RAISE EXCEPTION 'Invalid item_id: No matching record in dashboards for item_type = dashboard';
        END IF;
    ELSE
        RAISE EXCEPTION 'Invalid item_type: %', NEW.item_type;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER validate_contact_shares_item
BEFORE INSERT OR UPDATE ON contact_shares
FOR EACH ROW
EXECUTE FUNCTION validate_item_id();
