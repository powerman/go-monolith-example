-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE PROCEDURE down_not_supported()
    SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT='downgrade is not supported, restore from backup instead';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
-- Example usage: CALL down_not_supported();
DROP PROCEDURE down_not_supported;
