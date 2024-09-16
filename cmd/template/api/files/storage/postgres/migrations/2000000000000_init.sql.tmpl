-- +goose Up
-- +goose StatementBegin
CREATE FUNCTION sync_updated_at_column () RETURNS TRIGGER LANGUAGE plpgsql AS $$
BEGIN
    IF NEW.deleted_at IS NULL THEN
        NEW.updated_at = CLOCK_TIMESTAMP();
    END IF;
    RETURN NEW;
END;
$$;

CREATE TABLE
    foos (
        id serial PRIMARY KEY,
        NAME VARCHAR(50) NOT NULL,
        created_at timestamptz DEFAULT CLOCK_TIMESTAMP(),
        updated_at timestamptz,
        deleted_at timestamptz
    );

CREATE TRIGGER sync_foo_updated_at BEFORE
UPDATE ON foos FOR EACH ROW
EXECUTE PROCEDURE sync_updated_at_column ();

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE foos;

DROP FUNCTION sync_updated_at_column;

-- +goose StatementEnd