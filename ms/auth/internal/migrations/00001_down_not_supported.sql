-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- Usage:
--      -- +goose Down
--      SELECT down_not_supported();
-- +goose StatementBegin
CREATE FUNCTION down_not_supported() RETURNS void LANGUAGE plpgsql AS $$
    BEGIN
        RAISE EXCEPTION 'downgrade is not supported, restore from backup instead';
    END;
$$;
-- +goose StatementEnd

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP FUNCTION down_not_supported;
