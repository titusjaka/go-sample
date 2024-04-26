-- +migrate Up
CREATE TABLE snippets
(
	id         serial                      NOT NULL PRIMARY KEY,
	title      text                        NOT NULL,
	content    text                        NOT NULL,
	created_at timestamp WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at timestamp WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
	expires_at timestamp WITHOUT TIME ZONE NOT NULL
);

CREATE INDEX idx_snippets_created_at ON snippets (created_at);

CREATE INDEX idx_snippets_expires_at ON snippets (expires_at);

-- +migrate Down
DROP TABLE snippets;
