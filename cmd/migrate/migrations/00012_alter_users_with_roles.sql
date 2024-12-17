-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS users
ADD COLUMN role_id INT REFERENCES roles (id) DEFAULT 1;

UPDATE users
SET
    role_id = (
        SELECT
            id
        FROM
            roles
        WHERE
            name = 'user'
    );

ALTER TABLE users
ALTER COLUMN role_id
DROP DEFAULT;

ALTER TABLE users
ALTER COLUMN role_id
SET
    NOT NULL;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS users
DROP COLUMN IF EXISTS role_id;

-- +goose StatementEnd