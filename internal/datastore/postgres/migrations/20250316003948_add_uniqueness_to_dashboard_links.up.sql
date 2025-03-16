ALTER TABLE dashboard_links ADD CONSTRAINT dashboard_links_contact_dashboard_unique UNIQUE (contact_id, dashboard_id);
