-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    IF NOT EXISTS roles (
        id BIGSERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL UNIQUE,
        level int NOT NULL DEFAULT 0,
        description TEXT
    );

INSERT INTO
    roles (name, description, level)
VALUES
    (
        'user',
        'A user can read and write to their own data',
        1
    );


INSERT INTO
    roles (name, description, level)
VALUES
    (
        'seller',
        'A seller can sell products',
        2
    );

INSERT INTO
    roles (name, description, level)
VALUES
    (
        'admin',
        'An admin can do everything',
        3
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS roles;

-- +goose StatementEnd