-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "user" (
	id bigint GENERATED ALWAYS AS IDENTITY NOT NULL,
	login varchar NOT NULL,
	password text NOT NULL,
	last_update timestamp without time zone DEFAULT now(),
	CONSTRAINT user_pk PRIMARY KEY (id),
	CONSTRAINT user_unique UNIQUE (login)
);
CREATE TABLE IF NOT EXISTS metadata (
	id bigint GENERATED ALWAYS AS IDENTITY NOT NULL,
	user_id bigint NOT NULL,
	title text NOT NULL,
	card_number text NULL,
	login text NULL,
	password text NULL,
	dtype t_dtype NOT NULL,

	CONSTRAINT metadata_user_fk_1 FOREIGN KEY (user_id) REFERENCES "user"(id) ON DELETE CASCADE,
	CONSTRAINT metadata_pk PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user;
DROP TABLE metadata;
-- +goose StatementEnd
