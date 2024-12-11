-- +goose Up
-- +goose StatementBegin
ALTER TABLE products
ADD COLUMN version INT NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE products
DROP COLUMN version;
-- +goose StatementEnd
