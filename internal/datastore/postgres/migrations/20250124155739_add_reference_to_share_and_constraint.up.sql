ALTER TABLE contact_shares ADD COLUMN item_reference VARCHAR(220) NOT NULL;

ALTER TABLE contact_shares ADD CONSTRAINT unique_contact_per_shared_item UNIQUE(item_reference,contact_id);
