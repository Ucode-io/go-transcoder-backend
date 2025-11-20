ALTER TABLE pipelines ADD webhook_status BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE pipelines ADD webhook_retry_count INT NOT NULL DEFAULT 0;
ALTER TABLE pipelines ADD webhook_last_retry TIMESTAMP;
