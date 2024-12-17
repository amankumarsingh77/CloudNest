ALTER TABLE user_quotas DROP CONSTRAINT positive_storage;
ALTER TABLE user_quotas
    ADD CONSTRAINT positive_storage CHECK (
        storage_used >= 0 AND
        storage_limit >= 0
        );
