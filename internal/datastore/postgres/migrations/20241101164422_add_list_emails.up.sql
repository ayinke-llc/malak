CREATE TABLE IF NOT EXISTS contact_list_mappings (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  list_id uuid NOT NULL REFERENCES contact_lists(id),
  contact_id uuid NOT NULL REFERENCES contacts(id),
  reference VARCHAR (220) UNIQUE NOT NULL,
  created_by uuid NOT NULL REFERENCES users(id),

  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE
);

ALTER TABLE contact_list_mappings ADD CONSTRAINT contact_list_emails_reference_check_key CHECK (reference ~ 'list_email_[a-zA-Z0-9._]+');
