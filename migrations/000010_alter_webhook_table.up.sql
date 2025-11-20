CREATE UNIQUE INDEX idx_unique_columns_project_id_webhook_type
ON webhooks (project_id, webhook_type);
