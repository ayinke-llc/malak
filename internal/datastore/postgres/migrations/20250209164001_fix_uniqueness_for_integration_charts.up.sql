ALTER TABLE integration_charts DROP CONSTRAINT unique_internal_name_per_workspace_and_integration;
ALTER TABLE integration_charts ADD CONSTRAINT unique_name_and_internal_name_per_workspace_and_integration UNIQUE(user_facing_name,internal_name,workspace_id,workspace_integration_id);
