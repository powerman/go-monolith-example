-- +goose Up

-- Usage:
--      CREATE TABLE {table} (
--          updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
--          ... other fields
--      );
--      CREATE TRIGGER {table}_updated_at
--          BEFORE UPDATE ON {table}
--          FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();
-- +goose StatementBegin
CREATE FUNCTION trigger_set_updated_at() RETURNS trigger LANGUAGE plpgsql AS $$
    BEGIN
        NEW.updated_at := now();
        RETURN NEW;
    END;
$$;
-- +goose StatementEnd

-- +goose Down
DROP FUNCTION trigger_set_updated_at;
