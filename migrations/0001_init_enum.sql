-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 't_dtype') THEN
        CREATE TYPE "t_dtype" AS ENUM ('CARD', 'LOGPASS', 'FILE');
    END IF;
END $$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TYPE IF EXISTS t_dtype;
-- +goose StatementEnd
