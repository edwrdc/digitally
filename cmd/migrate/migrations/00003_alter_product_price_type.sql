-- +goose Up
-- +goose StatementBegin
ALTER TABLE products
    ALTER COLUMN price TYPE NUMERIC(10,2) USING price::numeric;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE products
    ALTER COLUMN price TYPE TEXT USING price::text;
-- +goose StatementEnd 