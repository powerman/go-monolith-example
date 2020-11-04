-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE example MODIFY user_id VARCHAR(64) NOT NULL;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
CALL down_not_supported();
