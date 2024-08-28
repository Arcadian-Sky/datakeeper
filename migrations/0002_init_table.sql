-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "user" (
	id bigint GENERATED ALWAYS AS IDENTITY NOT NULL,
	login varchar NOT NULL,
	password text NOT NULL,
	CONSTRAINT user_pk PRIMARY KEY (id),
	CONSTRAINT user_unique UNIQUE (login)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user;
-- +goose StatementEnd
